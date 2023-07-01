package main

import (
	"fmt"
	"math/big"
	"os"
	"sync"
)

type AnagramFinder interface {
	FindAnagrams(wordsFile string) [][]string
}

type PrimeMultiplication struct {
	PrimeMap map[rune]int
	Settings AnagramSettings
}

func (solver PrimeMultiplication) FindAnagrams() [][]string {

}

func (solver PrimeMultiplication) ProcessLine(s string, wg *sync.WaitGroup) {
	defer wg.Done()

	primeLetters := solver.Settings.RegExp.ReplaceAllString(s, "")

	if len(primeLetters) <= solver.Settings.MaxLength {
		product := solver.MultiplyLetters(s)
		solver.AppendResultFile(s, product)
	}
}

func (solver PrimeMultiplication) MultiplyLetters(s string) *big.Int {
	product := big.NewInt(1)
	for _, char := range s {
		// Multiply letters by given PrimeMap value
		if char >= 'a' && char <= 'z' {
			factor := big.NewInt(int64(solver.PrimeMap[char]))
			product.Mul(product, factor)
		} else {
			product.Mul(product, big.NewInt(1))
		}
	}

	return product
}

func (solver PrimeMultiplication) AppendResultFile(anagram string, r *big.Int) {
	fileName := fmt.Sprintf("data/%s", r.String())
	f, err := os.OpenFile(fileName, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0600)
	if err != nil {
		// TODO change this
		panic(err)
	}
	defer f.Close()
	fmt.Fprintf(f, "%v\n", anagram)
}

//func Combinations(list []string, combinationLength int) []string {
//    if combinationLength > len(list) {
//		return []string{}
//	}
//	result := []string{}
//	combine(list, combinationLength, 0, []string{}, &result)
//	return result
//}
//
//func combine(list []string, combinationLength int, start int, currentCombination []string, result *[]string) {
//	if combinationLength == 0 {
//		*result = append(*result, strings.Join(currentCombination, " "))
//		return
//	}
//	for i := start; i <= len(list)-combinationLength; i++ {
//		combine(list, combinationLength-1, i+1, append(currentCombination, list[i]), result)
//	}
//}
//
