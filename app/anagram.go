package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"math/big"
	"os"
	"path/filepath"
	"sync"
)

type AnagramFinder interface {
	FindAnagrams(wordsFile string) ([][]string, error)
}

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

//
//type BitwiseMatch struct {
//	LetterFrequencyDescending []rune
//	Settings AnagramSettings
//}
//
//func (finder BitwiseMatch) FindAnagrams() ([][]string, error) {}
//func (finder BitwiseMatch) ProcessLine()
//func (finder BitwiseMatch) EncodeLettersToBitsAndWeights()
//func (finder BitwiseMatch) CombineAnagrams()

type PrimeMultiplication struct {
	BaseAlgorithm
	PrimeMap map[rune]int
}

func (finder PrimeMultiplication) FindAnagrams() ([][]string, error) {
	scanner, err := finder.BaseAlgorithm.ScanDataSource()
	if err != nil {
		return nil, err
	}

	var wg sync.WaitGroup
	for scanner.Scan() {
		wg.Add(1)
		go finder.ProcessLine(scanner.Text(), &wg)
	}
	wg.Wait()

	if err := scanner.Err(); err != nil {
		return nil, err
	}

	for i := 1; i < finder.BaseAlgorithm.Settings.MaxWords; i++ {
		go func() {
			finder.CombineResultFiles()
		}()
	}

	// placeholder
	res := [][]string{{}}
	return res, nil
}

func (finder PrimeMultiplication) ProcessLine(s string, wg *sync.WaitGroup) {
	defer wg.Done()

	primeLetters := finder.BaseAlgorithm.Settings.RegExp.ReplaceAllString(s, "")

	if len(primeLetters) <= finder.BaseAlgorithm.Settings.MaxLength {
		product := finder.MultiplyLetters(s)
		finder.BaseAlgorithm.AppendResultFile([]string{s}, product)
	}
}

func (finder PrimeMultiplication) MultiplyLetters(s string) *big.Int {
	product := big.NewInt(1)
	for _, char := range s {
		// Multiply letters by given PrimeMap value
		if char >= 'a' && char <= 'z' {
			factor := big.NewInt(int64(finder.PrimeMap[char]))
			product.Mul(product, factor)
		}
	}

	return product
}

func (finder PrimeMultiplication) CombineAnagrams(res1, res2 string) { //([]string, *big.Int) {
	product := big.NewInt(1)
	factor1, _ := new(big.Int).SetString(res1, 10)
	factor2, _ := new(big.Int).SetString(res2, 10)
	product.Mul(factor1, factor2)

	wd := finder.BaseAlgorithm.Settings.WorkDir()
	lines1, _ := ReadFile(filepath.Join(wd, res1))
	lines2, _ := ReadFile(filepath.Join(wd, res2))

	results := finder.BaseAlgorithm.CombineLines(lines1, lines2)
	finder.BaseAlgorithm.AppendResultFile(results, product)
}

func (finder PrimeMultiplication) CombineResultFiles() error {
	wd := finder.BaseAlgorithm.Settings.WorkDir()
	fs, err := ioutil.ReadDir(wd)
	if err != nil {
		return err
	}

	for _, f1 := range fs {
		for _, f2 := range fs[1:] {
			go finder.CombineAnagrams(f1.Name(), f2.Name())
		}

	}

	return err
}
