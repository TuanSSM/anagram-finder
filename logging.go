package main

import (
	"fmt"
	"time"
)

type LoggingService struct {
	next Service
}

func NewLoggingService(next Service) Service {
	return &LoggingService{
		next: next,
	}
}

func (s *LoggingService) GetAllDataSources() (res []DataSource, err error) {
	defer func(start time.Time) {
		fmt.Printf("INFO | response=%v\n", res)
	}(time.Now())

	return s.next.GetAllDataSources()
}

func (s *LoggingService) GetDataSource(uuid string) (res []string, err error) {
	defer func(start time.Time) {
		fmt.Printf("INFO | Datasource has %v lines\n", len(res))
	}(time.Now())

	return s.next.GetDataSource(uuid)
}

func (s *LoggingService) GrabDataSource(req *GrabDataSourceRequest) (ds *DataSource, err error) {
	defer func(start time.Time) {
		fmt.Printf("INFO | Datasource grabbed succesfully %v\n", ds)
	}(time.Now())

	return s.next.GrabDataSource(req)
}
