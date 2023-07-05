package main

import (
	"context"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/gofiber/fiber/v2/middleware/requestid"
	"github.com/gofiber/template/html/v2"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ApiServer struct {
}

func NewApiServer() *ApiServer {
	return &ApiServer{}
}

// Configure server
func (s *ApiServer) Start(listenAddr, mongoUri string) error {
	engine := html.New("web/template", ".tmpl")
	engine.Reload(true)
	app := fiber.New(fiber.Config{
		Views: engine,
	})

	app.Static("/", "web/static")
	app.Use(logger.New())
	app.Use(requestid.New())

	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI(mongoUri))
	if err != nil {
		panic(err)
	}

	store := client.Database("anagram-finder")
	dsStore := NewDatasourceStore(store)
	dsHandler := NewDatasourceHandler(dsStore)

	// Routes
	datasourceApp := app.Group("/datasource")
	datasourceApp.Get("", dsHandler.HandleGetDatasources)
	datasourceApp.Get(":id", dsHandler.HandleGetDatasourceByID)
	//datasourceApp.Post("", dsHandler.HandlePostDatasource)
	//datasourceApp.Get("/:uuid/content", s.handleGetDataSourceContent)
	//datasourceApp.Get("/:dictId/metrics", s.svc.getDictionaryMetrics)
	//app.Post("/solve", s.handleFindAnagrams)

	aStore := NewAnagramStore(store)
	aHandler := NewAnagramHandler(*aStore)

	app.Post("/find", aHandler.HandleFetchAnagramsFromUrl)

	err = app.Listen(listenAddr)

	return err
}
