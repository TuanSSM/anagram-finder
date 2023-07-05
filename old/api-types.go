package main

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
