package main

import (
	"bufio"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/cavaliergopher/grab/v3"
)

type Service interface {
	GetAllDataSources() ([]DataSource, error)
	GetDataSource(uuid string) ([]string, error)
	GrabDataSource(req *GrabDataSourceRequest) (*DataSource, error)
}

type AnagramService struct{}

func NewAnagramService() Service {
	return &AnagramService{}
}

func (s *AnagramService) GetAllDataSources() ([]DataSource, error) {
	dsFile := "data/datasources.csv"

	err := createFileIfDNE(dsFile)
	if err != nil {
		return nil, err
	}

	file, err := os.Open(dsFile)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var dataSources []DataSource
	for _, record := range records {
		dataSources = append(dataSources, DataSource{
			Name:   record[0],
			UUID:   record[1],
			RawUrl: record[2],
		})
	}

	return dataSources, nil
}

func (s *AnagramService) GetDataSource(uuid string) ([]string, error) {
	dataSources, err := s.GetAllDataSources()
	if err != nil {
		return nil, err
	}

	var filePath string
	for _, ds := range dataSources {
		if ds.UUID == uuid {
			filePath = ds.FilePath()
		}
	}
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	var lines []string
	scanner := bufio.NewScanner(f)

	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	return lines, nil
}

func (s *AnagramService) GrabDataSource(req *GrabDataSourceRequest) (*DataSource, error) {
	dataSources, err := s.GetAllDataSources()
	if err != nil {
		return nil, err
	}

	for _, ds := range dataSources {
		if ds.RawUrl == req.RawUrl {
			return nil, fmt.Errorf("file already downloaded")
		}
	}
	dataSource := NewDataSource(req.Name, req.RawUrl)

	resp, err := grab.Get(dataSource.FilePath(), req.RawUrl)
	fmt.Println(resp)
	if err != nil {
		return nil, err
	}

	dataSources = append(dataSources, *dataSource)

	file, err := os.OpenFile("data/datasources.csv", os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for _, ds := range dataSources {
		err := writer.Write([]string{ds.Name, ds.UUID, ds.RawUrl})
		if err != nil {
			return nil, err
		}
	}

	return dataSource, nil
}
