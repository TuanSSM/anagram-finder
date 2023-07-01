package main

import (
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

func (ds DataSource) ResultsPath() string {
	return "./app/data/" + ds.UUID + "/anagrams/"
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
	MaxWords     int    `json:"maxWords"`
	MaxLength    int    `json:"maxLength"`
}

type AnagramSettings struct {
	DataSource *DataSource    `json:"dataSource"`
	MaxWords   int            `json:"maxWords"`
	MaxLength  int            `json:"maxLength"`
	RegExp     *regexp.Regexp `json:"regExp"`
}

func NewAnagramSettings(ds *DataSource, mw int, ml int) *AnagramSettings {
	reg, err := regexp.Compile("[^a-z]+")
	if err != nil {
		return nil
	}

	return &AnagramSettings{
		DataSource: ds,
		MaxWords:   mw,
		MaxLength:  ml,
		RegExp:     reg,
	}
}
