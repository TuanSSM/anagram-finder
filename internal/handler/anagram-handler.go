package handler

import (
	"github.com/gofiber/fiber/v2"
	"github.com/tuanssm/anagram-finder/internal/manager"
	"github.com/tuanssm/anagram-finder/internal/store"
	"github.com/tuanssm/anagram-finder/internal/types"
)

type AnagramHandler struct {
	store store.AnagramStore
}

func NewAnagramHandler(aStore store.AnagramStore) *AnagramHandler {
	return &AnagramHandler{
		store: aStore,
	}
}

func (a *AnagramHandler) HandleCreateAnagramsFromUrl(c *fiber.Ctx) error {
	req := &types.FindAnagramsRequest{}
	if err := c.BodyParser(req); err != nil {
		return err
	}

	am := manager.NewAnagramManager(a.store, req.Datasource)
	err := am.ProcessURL(req.Datasource.RawUrl)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON("init ok")
}

//func (a *AnagramHandler) HandleGetAnagrams(c *fiber.Ctx) error {
//	phrase := c.Params("phrase")
//	res :=
//}
