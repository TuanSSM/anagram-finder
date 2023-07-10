package handler

import (
	"context"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/tuanssm/anagram-finder/internal/anagram"
	"github.com/tuanssm/anagram-finder/internal/store"
	"github.com/tuanssm/anagram-finder/internal/types"
)

type AnagramHandler struct {
	store  store.AnagramStore
	dstore store.DatasourceStorer
}

func NewAnagramHandler(aStore store.AnagramStore, dStore store.DatasourceStorer) *AnagramHandler {
	return &AnagramHandler{
		store:  aStore,
		dstore: dStore,
	}
}

func (a *AnagramHandler) HandleCreateAnagramsFromUrl(c *fiber.Ctx) error {
	req := &types.FindAnagramsRequest{}
	if err := c.BodyParser(req); err != nil {
		return err
	}

	ctx := context.Background()
	ds, err := a.dstore.GetByID(ctx, req.DatasourceId)
	if err != nil {
		return err
	}
	log.Printf("%v", ds)

	af := anagram.NewAnagramFinder(a.store, *ds)
	err = af.ProcessURL(ctx, req.MaxWords)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON("Anagrams are being created, this operation might take a while")
}

func (a *AnagramHandler) HandleGetAllAnagrams(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx := context.Background()
	ds, err := a.dstore.GetByID(ctx, id)
	if err != nil {
		return err
	}

	af := anagram.NewAnagramFinder(a.store, *ds)
	anagrams, err := af.GetAllAnagrams(ctx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(anagrams)
}
