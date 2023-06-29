package main

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
)

type ApiServer struct {
	svc Service
}

func NewApiServer(svc Service) *ApiServer {
	return &ApiServer{
		svc: svc,
	}
}

func (s *ApiServer) Start(listenAddr string) error {
	app := fiber.New()

	app.Use(logger.New())
	app.Use(requestid.New())

	app.Get("/", s.handleRoot)

	datasourceApp := app.Group("/datasource")
	datasourceApp.Get("/all", s.handleGetAllDataSources)
	datasourceApp.Get("/:uuid", s.handleGetDataSource)
	//datasourceApp.Get("/:dictId/metrics", s.svc.getDictionaryMetrics)
	datasourceApp.Post("", s.handleGrabDataSource)
	err := app.Listen(listenAddr)

	return err
}

func (s *ApiServer) handleRoot(ctx *fiber.Ctx) error {
	return ctx.Status(fiber.StatusOK).SendString("Anagram Finder")
}

func (s *ApiServer) handleGetAllDataSources(ctx *fiber.Ctx) error {
	res, err := s.svc.GetAllDataSources()
	if err != nil {
		return ctx.Status(fiber.StatusInternalServerError).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (s *ApiServer) handleGetDataSource(ctx *fiber.Ctx) error {
	res, err := s.svc.GetDataSource(ctx.Params("uuid"))
	if err != nil {
		return ctx.Status(fiber.StatusNotFound).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}

func (s *ApiServer) handleGrabDataSource(ctx *fiber.Ctx) error {
	req := new(GrabDataSourceRequest)
	if err := ctx.BodyParser(req); err != nil {
		return ctx.Status(fiber.StatusBadRequest).SendString(err.Error())
	}

	res, err := s.svc.GrabDataSource(req)
	if err != nil {
		return ctx.Status(fiber.StatusUnprocessableEntity).SendString(err.Error())
	}

	return ctx.Status(fiber.StatusOK).JSON(res)
}
