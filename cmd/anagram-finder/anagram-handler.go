package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
)

type AnagramHandler struct {
	store AnagramStore
}

func NewAnagramHandler(aStore AnagramStore) *AnagramHandler {
	return &AnagramHandler{
		store: aStore,
	}
}

func (a *AnagramHandler) HandleFetchAnagramsFromUrl(c *fiber.Ctx) error {
	req := &FindAnagramsRequest{}
	if err := c.BodyParser(req); err != nil {
		return err
	}

	as := NewAnagramSearch(context.Background(), a.store, req.Datasource)
	err := as.ProcessURL(req.Datasource.RawUrl)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON("init ok")
}

//
//	//datasource, err := NewFindAnagramsRequest(req)
//	//if err != nil {
//	//	return err
//	//}
//
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	if err := h.store.Insert(ctx, datasource); err != nil {
//		return err
//	}
//
//
//	//FetchInsertAnagrams(datasource.RawUrl)
//	//resp, err := grab.Get(datasource.Slug, datasource.RawUrl)
//	//if err != nil {
//	//	return err
//	//}
//	//log.Printf("File %s is downloaded", resp.Filename)
//
//	return c.Status(fiber.StatusOK).JSON(datasource)
