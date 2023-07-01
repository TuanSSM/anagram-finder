package main

import (
	"fmt"
	"regexp"

	"github.com/google/uuid"
)

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
	return "./app/data/" + ds.UUID + "/" + ds.Name
}

func (ds DataSource) DirName() string {
	return "./app/data/" + ds.UUID
}

type DataSourcesResponse struct {
	Datasources []DataSource `json:"datasources"`
}

type GrabDataSourceRequest struct {
	RawUrl string `json:"rawUrl"`
	Name   string `json:"name"`
}

type FindAnagramsRequest struct {
	DictionaryId string `json:"dictionaryId"`
	Algo         string `json:"algo"`
	MaxWords     int    `json:"maxWords"`
	MaxLength    int    `json:"maxLength"`
}

type AnagramSettings struct {
	DataSource *DataSource    `json:"dataSource"`
	Algo       string         `json:"algo"`
	MaxWords   int            `json:"maxWords"`
	MaxLength  int            `json:"maxLength"`
	RegExp     *regexp.Regexp `json:"regExp"`
}

func NewAnagramSettings(ds *DataSource, algo string, mw int, ml int) *AnagramSettings {
	reg, err := regexp.Compile("[^a-z]+")
	if err != nil {
		return nil
	}

	return &AnagramSettings{
		DataSource: ds,
		Algo:       algo,
		MaxWords:   mw,
		MaxLength:  ml,
		RegExp:     reg,
	}
}

func (as AnagramSettings) WorkDir() string {
	wd := fmt.Sprintf("/work/%s_%d_%d", as.Algo, as.MaxWords, as.MaxLength)
	return as.DataSource.DirName() + wd
}
