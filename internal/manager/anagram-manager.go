package manager

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/tuanssm/anagram-finder/internal/store"
	"github.com/tuanssm/anagram-finder/internal/types"
	"github.com/tuanssm/anagram-finder/internal/util"
)

type AnagramManager struct {
	store store.AnagramStore
	ds    types.Datasource
}

func NewAnagramManager(store store.AnagramStore, ds types.Datasource) *AnagramManager {
	return &AnagramManager{
		store: store,
		ds:    ds,
	}
}

func (as *AnagramManager) ProcessURL(url string) error {
	ctx := context.Background()

	chunkSize := 1000
	datach := make(chan string)
	anagramch := make(chan *types.AnagramEntry)
	errCh := make(chan error, 3)
	log.Println("[ INFO ] Channels created.")

	go func() {
		log.Println("[ INFO ] Fetching Lines.")
		if err := FetchLines(ctx, url, datach); err != nil {
			errCh <- err
		}
	}()

	go func() {
		log.Println("[ INFO ] Parsing Lines.")
		if err := ParseLines(ctx, datach, anagramch); err != nil {
			errCh <- err
			return
		}
	}()

	log.Printf("[ INFO ] Inserting Anagrams to %s", as.ds.Slug)
	go InsertManyEntries(ctx, &as.store, as.ds.Slug, anagramch, chunkSize)

	for i := 0; i < 3; i++ {
		if err := <-errCh; err != nil {
			return err
		}
	}

	return nil
}

func FetchLines(ctx context.Context, url string, linech chan<- string) error {
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %s", util.ErrFailedToFetch, err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case linech <- scanner.Text():
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("%w: %s", util.ErrFailedToFetch, scanner.Err())
	}

	return nil
}

func ParseLines(ctx context.Context, linech <-chan string, anagramch chan<- *types.AnagramEntry) error {
	i := 0
	for line := range linech {
		i++
		if i%1000 == 0 {
			log.Printf("Processed Line #%d: %s", i, line)
		}
		anagram := types.NewAnagramEntry(line)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case anagramch <- anagram:
		}
	}

	return nil
}

func InsertManyEntries(ctx context.Context, store *store.AnagramStore, coll string, anagramch <-chan *types.AnagramEntry, chunkSize int) error {
	var chunk []*types.AnagramEntry

	i := 0
	for anagram := range anagramch {
		chunk = append(chunk, anagram)

		if len(chunk) >= chunkSize {
			err := store.BulkInsert(ctx, coll, chunk)
			if err != nil {
				return fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
			}
			log.Printf("Inserted chunk %d-%d to anagrams", (i*chunkSize + 1), (i+1)*chunkSize)
			i++
			chunk = nil
		}
	}

	if len(chunk) > 0 {
		err := store.BulkInsert(ctx, coll, chunk)
		if err != nil {
			return fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
		}
		log.Printf("Inserted remaining %d to anagrams", len(chunk))
	}

	return nil
}

func InsertEntries(ctx context.Context, store *store.AnagramStore, coll string, anagramch <-chan *types.AnagramEntry) error {
	for anagram := range anagramch {
		log.Printf("Anagram: %v", anagram)
		err := store.Insert(ctx, coll, anagram)
		if err != nil {
			return fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
		}
	}

	return nil
}
