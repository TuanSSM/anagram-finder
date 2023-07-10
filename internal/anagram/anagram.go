package anagram

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/tuanssm/anagram-finder/internal/store"
	"github.com/tuanssm/anagram-finder/internal/types"
	"github.com/tuanssm/anagram-finder/internal/util"
)

type AnagramFinder struct {
	store  store.AnagramStore
	ds     types.Datasource
	angMap map[types.BitWeights][]string
}

func NewAnagramFinder(store store.AnagramStore, ds types.Datasource) *AnagramFinder {
	angMap := make(map[types.BitWeights][]string)
	return &AnagramFinder{
		store:  store,
		ds:     ds,
		angMap: angMap,
	}
}

func (af *AnagramFinder) ProcessURL(ctx context.Context, maxWords int) error {
	errCh := make(chan error, 6)

	err := os.MkdirAll(af.ds.WorkDir(), os.ModePerm)
	if err != nil {
		log.Println(err)
	}

	dataCh := af.FetchLines(ctx, errCh)
	encodeCh := af.ParseLines(ctx, dataCh, errCh)
	anagrams := af.MatchBaseAnagrams(encodeCh)
	combCh := af.GenerateCombinations(anagrams, maxWords)

	go af.AppendAnagramsToJSON(ctx, combCh, errCh)

	batchSize := 100
	namesCh := af.ReadDirBatches(batchSize, errCh)
	anagCh := af.AddAnagramsToChannel(ctx, namesCh, errCh)
	go af.BulkInsertAnagrams(ctx, 10000, anagCh, errCh)

	for err := range errCh {
		if err != nil {
			return err
		}
	}

	return nil
}

func (af *AnagramFinder) ShowAnagrams(ctx context.Context) ([]*types.AnagramEntry, error) {
	anagrams, err := af.GetAllAnagrams(ctx)
	if err != nil {
		return nil, err
	}
	return anagrams, nil
}

func (af *AnagramFinder) GenerateCombinations(angs []*types.AnagramEntry, max int) <-chan *types.AnagramEntry {
	combCh := make(chan *types.AnagramEntry, 1000)

	log.Printf("Generating combinations from %d items", len(angs))
	go func(c chan *types.AnagramEntry) {
		defer close(c)
		baseEnc := &types.BitWeights{}
		baseAnagram := &types.AnagramEntry{
			Encoded:  *baseEnc,
			Anagrams: []string{""},
		}
		af.CombineAnagrams(c, baseAnagram, angs, max)
	}(combCh)

	return combCh
}

func (af *AnagramFinder) CombineAnagrams(c chan *types.AnagramEntry, combo *types.AnagramEntry, angs []*types.AnagramEntry, length int) {
	if length <= 0 {
		return
	}

	var newCombo *types.AnagramEntry
	for _, ch := range angs {
		newCombo = combo.Combine(ch)
		c <- newCombo
		af.CombineAnagrams(c, newCombo, angs, length-1)
	}
}

func (af *AnagramFinder) MatchBaseAnagrams(angCh <-chan *types.AnagramEntry) []*types.AnagramEntry {
	var anagrams []*types.AnagramEntry
	i := 1
	for ang := range angCh {
		if _, ok := af.angMap[ang.Encoded]; ok {
			log.Printf("[ INFO ] Found anagrams: %v, %v", af.angMap[ang.Encoded], ang.Anagrams)
			af.angMap[ang.Encoded] = append(af.angMap[ang.Encoded], ang.Anagrams...)
		} else {
			af.angMap[ang.Encoded] = ang.Anagrams
		}
		if i%1000 == 0 {
			log.Printf("[ INFO ] Encodings processed: %d", i)
		}
		i++
	}
	log.Printf("[ INFO ] Total encodings indexed is: %d", len(af.angMap))

	for bw, strs := range af.angMap {
		anagrams = append(anagrams, &types.AnagramEntry{
			Encoded:  bw,
			Anagrams: strs,
		})
	}

	return anagrams
}

func (a *AnagramFinder) AddAnagramsToChannel(ctx context.Context, fileCh <-chan []string, errCh chan<- error) <-chan *types.AnagramEntry {
	anagCh := make(chan *types.AnagramEntry, 10000)

	go func(anagCh chan *types.AnagramEntry, errCh chan<- error) {
		defer close(anagCh)
		for fs := range fileCh {
			for _, f := range fs {

				path := fmt.Sprintf("%s/%s", a.ds.WorkDir(), f) //filepath.Join(dir, f)
				log.Printf("%v", path)

				anagrams, err := types.ReadAnagramsFromJSON(path)
				if err != nil {
					errCh <- err
				}

				for _, a := range anagrams {
					select {
					case <-ctx.Done():
						errCh <- ctx.Err()
					case anagCh <- a:
					}

				}
			}
		}
	}(anagCh, errCh)

	return anagCh
}

func (a *AnagramFinder) ReadDirBatches(batchSize int, errCh chan<- error) <-chan []string {
	path := fmt.Sprintf("%s/", a.ds.WorkDir())
	dir, err := os.Open(path)
	if err != nil {
		errCh <- err
		return nil
	}

	namesCh := make(chan []string)

	go func(c chan []string) {
		defer dir.Close()
		defer close(c)
		names := make([]string, batchSize)
		for len(names) == batchSize {
			names, err = dir.Readdirnames(batchSize)
			if err != nil {
				fmt.Printf("error reading directory names: %v\n", err)
				errCh <- err
			}
			namesCh <- names
		}
	}(namesCh)

	return namesCh
}

func (af AnagramFinder) FetchLines(ctx context.Context, errCh chan<- error) <-chan string {
	log.Printf("[ INFO ] Fetching Lines from url %s.", af.ds.RawUrl)
	lineCh := make(chan string)
	go func(c chan string) {
		resp, err := http.Get(af.ds.RawUrl)
		if err != nil {
			errCh <- fmt.Errorf("%w: %s", util.ErrFailedToFetch, err)
		}

		defer resp.Body.Close()

		scanner := bufio.NewScanner(resp.Body)
		for scanner.Scan() {
			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
			case lineCh <- scanner.Text():
			}
		}
		close(lineCh)

		if scanner.Err() != nil {
			errCh <- fmt.Errorf("%w: %s", util.ErrFailedToFetch, scanner.Err())
		}
	}(lineCh)

	return lineCh
}

func (af AnagramFinder) ParseLines(ctx context.Context, linech <-chan string, errCh chan<- error) <-chan *types.AnagramEntry {
	anagramCh := make(chan *types.AnagramEntry)
	go func(c chan *types.AnagramEntry) {
		defer close(c)
		i := 0
		for line := range linech {
			i++
			if i%1000 == 0 {
				log.Printf("Processed Line #%d: %s", i, line)
			}
			anagram := types.NewAnagramEntry(line)

			select {
			case <-ctx.Done():
				errCh <- ctx.Err()
			case anagramCh <- anagram:
			}
		}
	}(anagramCh)

	return anagramCh
}

func (a AnagramFinder) AppendAnagramsToJSON(ctx context.Context, anagramch <-chan *types.AnagramEntry, errCh chan<- error) {
	i := 0
	log.Printf("Starting Append Worker")
	for anagram := range anagramch {
		i++
		err := anagram.AppendToJSON(a.ds.WorkDir())
		if err != nil {
			errCh <- err
		}
		if i%1000 == 0 {
			log.Printf("Appended anagram #%d: %v", i, anagram)
		}
	}
}

func (af AnagramFinder) BulkInsertAnagrams(ctx context.Context, chunkSize int, anagramch <-chan *types.AnagramEntry, errCh chan<- error) {
	defer close(errCh)
	var chunk []*types.AnagramEntry

	i := 0
	for anagram := range anagramch {
		chunk = append(chunk, anagram)

		if len(chunk) >= chunkSize {
			err := af.store.BulkInsert(ctx, af.ds.Slug, chunk)
			if err != nil {
				errCh <- fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
			}
			log.Printf("Inserted chunk %d-%d to anagrams", (i*chunkSize + 1), (i+1)*chunkSize)
			i++
			chunk = nil
		}
	}

	if len(chunk) > 0 {
		err := af.store.BulkInsert(ctx, af.ds.Slug, chunk)
		if err != nil {
			errCh <- fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
		}
		log.Printf("Inserted remaining %d to anagrams", len(chunk))
	}
}

//func (af AnagramFinder) InsertEntries(ctx context.Context, coll string, anagramch <-chan *types.AnagramEntry) error {
//	for anagram := range anagramch {
//		log.Printf("Anagram: %v", anagram)
//		err := af.store.Insert(ctx, coll, anagram)
//		if err != nil {
//			return fmt.Errorf("%w: %s", util.ErrFailedToInsert, err)
//		}
//	}
//
//	return nil
//}

func (af AnagramFinder) GetAllAnagrams(ctx context.Context) ([]*types.AnagramEntry, error) {
	anagrams, err := af.store.GetAll(ctx, af.ds.Slug)
	if err != nil {
		return nil, err
	}
	return anagrams, err
}
