package main

import (
	"bufio"
	"context"
	"fmt"
	"net/http"
)

type AnagramSearch struct {
	ctx   context.Context
	store AnagramStore
	ds    Datasource
}

func NewAnagramSearch(ctx context.Context, store AnagramStore, ds Datasource) *AnagramSearch {
	return &AnagramSearch{
		ctx:   ctx,
		store: store,
		ds:    ds,
	}
}

func (as *AnagramSearch) ProcessURL(url string) error {
	lines := make(chan string)
	entries := make(chan *AnagramEntry)
	errChan := make(chan error)

	ctx, cancel := context.WithCancel(as.ctx)
	defer cancel() // ensure all paths cancel the context to avoid context leak

	go func() {
		defer close(lines)
		if err := FetchLines(ctx, url, lines); err != nil {
			errChan <- err
			cancel()
		}
	}()

	go func() {
		defer close(entries)
		if err := ParseLines(ctx, lines, entries); err != nil {
			errChan <- err
			cancel()
		}
	}()

	go func() {
		if err := InsertEntries(ctx, &as.store, as.ds.Slug, entries); err != nil {
			errChan <- err
			cancel()
		}
	}()

	select {
	case <-ctx.Done():
		return ctx.Err() // context was canceled, return the reason
	case err := <-errChan:
		return err // received an error from one of the goroutines
	}
}

func FetchLines(ctx context.Context, url string, lines chan<- string) error {
	defer close(lines)

	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrFailedToFetch, err)
	}
	defer resp.Body.Close()

	scanner := bufio.NewScanner(resp.Body)
	for scanner.Scan() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case lines <- scanner.Text():
		}
	}

	if scanner.Err() != nil {
		return fmt.Errorf("%w: %s", ErrFailedToFetch, scanner.Err())
	}

	return nil
}

func ParseLines(ctx context.Context, lines <-chan string, entries chan<- *AnagramEntry) error {
	defer close(entries)

	for line := range lines {
		entry := NewAnagramEntry(line)

		select {
		case <-ctx.Done():
			return ctx.Err()
		case entries <- entry:
		}
	}

	return nil
}

func InsertEntries(ctx context.Context, store *AnagramStore, coll string, entries <-chan *AnagramEntry) error {
	for entry := range entries {
		err := store.Insert(ctx, coll, entry)
		if err != nil {
			return fmt.Errorf("%w: %s", ErrFailedToInsert, err)
		}
	}

	return nil
}
