package main

import "github.com/google/uuid"

type DataSource struct {
	Name   string `json:"name"`
	UUID   string `json:"uuid"`
	RawUrl string `json:"rawUrl"`
}

func NewDataSource(name string, rawUrl string) *DataSource {
	uuid := uuid.New().String()[:8]
	return &DataSource{
		Name:   name,
		UUID:   uuid,
		RawUrl: rawUrl,
	}
}

func (ds DataSource) FilePath() string {
	return "data/" + ds.UUID + "/" + ds.Name
}

type GrabDataSourceRequest struct {
	RawUrl string `json:"rawUrl"`
	Name   string `json:"name"`
}
