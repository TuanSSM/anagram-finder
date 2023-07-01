package main

import "errors"

var ErrDataSourceNotFound = errors.New("Datasource not found in datasources.csv")

var ErrDataSourceFileAccess = errors.New("Datasource file could not be opened")
