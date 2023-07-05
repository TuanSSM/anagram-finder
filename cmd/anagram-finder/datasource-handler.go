package main

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
)

type DatasourceHandler struct {
	store DatasourceStorer
}

func NewDatasourceHandler(dsStore DatasourceStorer) *DatasourceHandler {
	return &DatasourceHandler{
		store: dsStore,
	}
}

//func (h *DatasourceHandler) HandlePostDatasource(c *fiber.Ctx) error {
//	req := &CreateDatasourceRequest{}
//	if err := c.BodyParser(req); err != nil {
//		return err
//	}
//
//	datasource, err := NewDatasourceFromRequest(req.Datasources)
//	if err != nil {
//		return err
//	}
//
//	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
//	if err := h.store.Insert(ctx, datasource); err != nil {
//		return err
//	}
//
//	FetchInsertAnagrams(datasource.RawUrl)
//	//resp, err := grab.Get(datasource.Slug, datasource.RawUrl)
//	//if err != nil {
//	//	return err
//	//}
//	//log.Printf("File %s is downloaded", resp.Filename)
//
//	return c.Status(fiber.StatusOK).JSON(datasource)
//}

func (h *DatasourceHandler) HandleGetDatasources(c *fiber.Ctx) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	datasources, err := h.store.GetAll(ctx)
	if err != nil {
		return err
	}
	return c.Render("index", datasources)
}

func (h *DatasourceHandler) HandleGetDatasourceByID(c *fiber.Ctx) error {
	id := c.Params("id")
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	datasource, err := h.store.GetByID(ctx, id)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(datasource)
}
