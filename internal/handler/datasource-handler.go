package handler

import (
	"context"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/tuanssm/anagram-finder/internal/store"
	"github.com/tuanssm/anagram-finder/internal/types"
)

type DatasourceHandler struct {
	store store.DatasourceStorer
}

func NewDatasourceHandler(dsStore store.DatasourceStorer) *DatasourceHandler {
	return &DatasourceHandler{
		store: dsStore,
	}
}

func (h *DatasourceHandler) HandleGetDatasources(c *fiber.Ctx) error {
	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	datasources, err := h.store.GetAll(ctx)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(datasources)
	//return c.Render("index", datasources)
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

func (h *DatasourceHandler) HandlePostDatasource(c *fiber.Ctx) error {
	req := &types.CreateDatasourceRequest{}
	if err := c.BodyParser(req); err != nil {
		return err
	}

	ds, err := types.NewDatasourceFromRequest(req)
	if err != nil {
		return err
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)

	newDs, err := h.store.Insert(ctx, ds)
	if err != nil {
		return err
	}

	return c.Status(fiber.StatusOK).JSON(newDs)
}
