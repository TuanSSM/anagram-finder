package main

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
)

func createFileIfDNE(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		_, err := os.Create(path)
		if err != nil {
			return err
		}
	}

	return nil
}

func ReadFile(path string) ([]string, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)

	var lines []string
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	return lines, nil
}

func AppendFile(path string, lines []string) {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		if err := os.MkdirAll(filepath.Dir(path), 0770); err != nil {
			// TODO add custom err
			panic(err)
		}
		AppendFile(path, lines)
	}
	defer f.Close()
	for _, line := range lines {
		fmt.Fprintf(f, "%v\n", line)
	}
}
