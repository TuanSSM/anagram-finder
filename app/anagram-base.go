package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
)

type BaseAlgorithm struct {
	LetterFrequencyDescending []rune
	Settings                  AnagramSettings
}

func (finder BaseAlgorithm) ScanDataSource() (*bufio.Scanner, error) {
	fPath := finder.Settings.DataSource.FilePath()
	file, err := os.Open(fPath)
	if err != nil {
		return nil, ErrDataSourceFileAccess
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	return scanner, nil
}

func (finder BaseAlgorithm) CombineLines(lines1, lines2 []string) []string {
	var combinations []string
	for _, line1 := range lines1 {
		for _, line2 := range lines2 {
			combinations = append(combinations, line1+" "+line2)
		}
	}
	return combinations
}

func (finder BaseAlgorithm) AppendResultFile(anagrams []string, r *big.Int) {
	nComb := len(anagrams[0])
	var fPath string
	// Seperate anagrams by number of words contained
	if nComb == 1 {
		fPath = fmt.Sprintf("%s/%s", finder.Settings.WorkDir(), r.String())
	} else {
		fPath = fmt.Sprintf("%s/%d/%s", finder.Settings.WorkDir(), nComb, r.String())
	}

	AppendFile(fPath, anagrams)
}

func (finder BaseAlgorithm) SquashResultFiles() ([][]string, error) {
	wd := finder.Settings.WorkDir()
	filepath.Walk(wd, func(path string, info os.FileInfo, err error) error {
		//
		if info.IsDir() {
			filepath.Walk(filepath.Join(wd, info.Name()), func(path string, childInfo os.FileInfo, err error) error {

				return nil
			})
			return nil
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		newPath := filepath.Join(wd, info.Name())
		err = ioutil.WriteFile(newPath, data, info.Mode())
		if err != nil {
			return err
		}

		return nil
	})
}
