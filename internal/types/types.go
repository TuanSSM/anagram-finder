package types

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

type Datasource struct {
	ID     string `bson:"_id"`
	Name   string `bson:"name"`
	Slug   string `bson:"slug"`
	RawUrl string `bson:"rawUrl"`
}

type CreateDatasourceRequest struct {
	Name   string `bson:"name"`
	RawUrl string `bson:"rawUrl"`
}

func NewDatasourceFromRequest(req *CreateDatasourceRequest) (*Datasource, error) {
	//if err := validateCreateDatasourceRequest(req); err != nil {
	//	return nil, err
	//}

	parts := strings.Split(strings.ToLower(req.Name), " ")
	slug := strings.Join(parts, "-")

	return &Datasource{
		Name:   req.Name,
		Slug:   slug,
		RawUrl: req.RawUrl,
	}, nil
}

func (ds *Datasource) WorkDir() string {
	return fmt.Sprintf("./data/%s-%s", ds.Slug, ds.ID)
}

type FindAnagramsRequest struct {
	DatasourceId string `bson:"datasourceId"`
	MaxWords     int    `bson:"maxWords"`
}

func NewAnagramSettings(req *FindAnagramsRequest) {

}

type BitWeights struct {
	Bits    uint32  `bson:"encodedBits"`
	Weights [26]int `bson:"weights"`
}

func NewBitWeights(s string) *BitWeights {
	bw := &BitWeights{}
	alphabet := "esiarntolcdugpmhbyfvkwzxjq"
	for i, char := range alphabet {
		if strings.ContainsRune(s, char) {
			bw.Bits |= 1 << i
			bw.Weights[i] = strings.Count(s, string(char))
		}
	}

	return bw
}

func (bw *BitWeights) WeightsStr() string {
	str := make([]string, len(bw.Weights))
	for i, w := range bw.Weights {
		str[i] = strconv.Itoa(w)
	}
	return strings.Join(str, "-")
}

func WeightsFromStr(s string) ([26]int, error) {
	weights := strings.Split(s, "-")
	ws := [26]int{}
	for i, wstr := range weights {
		w, err := strconv.Atoi(wstr)
		if err != nil {
			return [26]int{}, err
		}
		ws[i] = int(w)
	}
	return ws, nil
}

func BitWeightsFromFile(s string) *BitWeights {
	parts := strings.Split(s, "-")

	first, err := strconv.ParseInt(parts[0], 10, 32)
	if err != nil {
		return nil
	}
	encodedBits := uint32(first)

	var weights [26]int
	for i := 1; i < len(parts); i++ {
		num, err := strconv.ParseInt(parts[i], 10, 32)
		if err != nil {
			return nil
		}
		weights[i-1] = int(num)
	}

	return &BitWeights{
		Bits:    encodedBits,
		Weights: weights,
	}
}

func (bw *BitWeights) Combine(other *BitWeights) *BitWeights {
	newBW := &BitWeights{
		Bits: bw.Bits | other.Bits,
	}

	for i := 0; i < 26; i++ {
		newBW.Weights[i] = bw.Weights[i] + other.Weights[i]
	}

	return newBW
}

func (bw *BitWeights) IsEquivalent(other *BitWeights) bool {
	return bw.Bits == other.Bits && bw.Weights == other.Weights
}

func (bw *BitWeights) ToString() string {
	str := make([]string, len(bw.Weights))
	for i, w := range bw.Weights {
		str[i] = strconv.Itoa(w)
	}
	wstr := strings.Join(str, "-")
	return fmt.Sprintf("%d-%s", bw.Bits, wstr)
}

type AnagramEntry struct {
	ID       string     `bson:"uuid"`
	Encoded  BitWeights `bson:"encoded"`
	Anagrams []string   `bson:"anagrams"`
}

func NewAnagramEntry(s string) *AnagramEntry {
	bw := NewBitWeights(s)
	return &AnagramEntry{
		Encoded:  *bw,
		Anagrams: []string{s},
	}
}

func (a *AnagramEntry) Combine(other *AnagramEntry) *AnagramEntry {
	var news []string
	bw := a.Encoded.Combine(&other.Encoded)
	for _, a1 := range a.Anagrams {
		var joinStr string
		if a1 == "" {
			joinStr = ""
		} else {
			joinStr = " "
		}

		var variations []string
		for _, a2 := range other.Anagrams {
			variation := strings.Join([]string{a1, a2}, joinStr)
			variations = append(variations, variation)
		}
		news = append(news, strings.Join(variations, ";"))
	}
	return &AnagramEntry{
		Encoded:  *bw,
		Anagrams: news,
	}
}

//func (a *AnagramEntry) AppendToFile(p string, numWords int) error {
//	fPath := fmt.Sprintf("%s/%d/%s", p, numWords, a.Encoded.ToString())
//	err := util.AppendLinesToFile(fPath, a.Anagrams)
//	if err != nil {
//		return err
//	}
//
//	return nil
//}

func (a *AnagramEntry) AppendToJSON(dir string) error {
	fName := fmt.Sprintf("%s/%d.json", dir, a.Encoded.Bits)

	fContent := make(map[string][]string)

	_, err := os.Stat(fName)
	if err == nil {
		fileContent, readErr := ioutil.ReadFile(fName)
		if readErr != nil {
			return readErr
		}

		jsonErr := json.Unmarshal(fileContent, &fContent)
		if jsonErr != nil {
			return jsonErr
		}
	}

	w := a.Encoded.WeightsStr()
	fContent[w] = append(fContent[w], a.Anagrams...)

	new, jsonErr := json.Marshal(fContent)
	if jsonErr != nil {
		return jsonErr
	}

	writeErr := ioutil.WriteFile(fName, new, 0644)
	if writeErr != nil {
		return writeErr
	}

	return nil
}

func ReadAnagramsFromJSON(path string) ([]*AnagramEntry, error) {
	var anagrams []*AnagramEntry

	fName := strings.TrimSuffix(filepath.Base(path), ".json")
	bits, parseErr := strconv.ParseUint(fName, 10, 32)
	if parseErr != nil {
		return nil, parseErr
	}
	fContent := make(map[string][]string)

	fileContent, readErr := ioutil.ReadFile(path)
	if readErr != nil {
		return nil, readErr
	}

	jsonErr := json.Unmarshal(fileContent, &fContent)
	if jsonErr != nil {
		return nil, jsonErr
	}

	for ws, as := range fContent {
		if len(as) < 2 {
			continue
		}
		w, err := WeightsFromStr(ws)
		if err != nil {
			return nil, err
		}
		bw := &BitWeights{
			Bits:    uint32(bits),
			Weights: w,
		}
		a := &AnagramEntry{
			Encoded:  *bw,
			Anagrams: as,
		}
		anagrams = append(anagrams, a)
	}
	return anagrams, nil
}
