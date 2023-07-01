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
	GetDataSource(uuid string) (*DataSource, error)
	GetDataSourceContent(uuid string) ([]string, error)
	GrabDataSource(req *GrabDataSourceRequest) (*DataSource, error)
	FindAnagrams(req *FindAnagramsRequest) ([][]string, error)
}

type AnagramService struct{}

func NewAnagramService() Service {
	return &AnagramService{}
}

func (s *AnagramService) GetAllDataSources() ([]DataSource, error) {
	dsFile := "./app/data/datasources.csv"

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

func (s *AnagramService) GetDataSource(uuid string) (*DataSource, error) {
	dataSources, err := s.GetAllDataSources()
	if err != nil {
		return nil, err
	}

	for _, ds := range dataSources {
		if ds.UUID == uuid {
			return &ds, nil
		}
	}

	return nil, ErrNotFound
}

func (s *AnagramService) GetDataSourceContent(uuid string) ([]string, error) {
	ds, err := s.GetDataSource(uuid)
	fPath := ds.FilePath()
	f, err := os.Open(fPath)
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

	file, err := os.OpenFile("./app/data/datasources.csv", os.O_WRONLY|os.O_CREATE, 0666)
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

func (s *AnagramService) FindAnagrams(req *FindAnagramsRequest) ([][]string, error) {
	ds, err := s.GetDataSource(req.DictionaryId)
	if err != nil {
		return nil, err
	}

	as := NewAnagramSettings(ds, req.MaxWords, req.MaxLength)

	primeSolver := PrimeMultiplication{
		PrimeMap: map[rune]int{
			'e': 2, 's': 3, 'i': 5, 'a': 7, 'r': 11, 'n': 13,
			't': 17, 'o': 19, 'l': 23, 'c': 29, 'd': 31,
			'u': 37, 'g': 41, 'p': 43, 'm': 47, 'h': 53,
			'b': 59, 'y': 61, 'f': 67, 'v': 71, 'k': 73,
			'w': 79, 'z': 83, 'x': 89, 'j': 97, 'q': 101,
		},
		Settings: *as,
	}

	res := primeSolver.FindAnagrams()
	return res, nil
}
