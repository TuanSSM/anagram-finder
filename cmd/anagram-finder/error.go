package main

import "errors"

var ErrDataSourceNotFound = errors.New("Datasource not found in datasources.csv")

var MongoNewClientError = errors.New("Mongo new client initialization has failed")

var ErrFailedToFetch = errors.New("Failed to fetch data from URL")
var ErrFailedToParse = errors.New("Failed to parse line into AnagramEntry")
var ErrFailedToInsert = errors.New("Failed to insert AnagramEntry into database")
