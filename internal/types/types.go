package types

import (
	"strings"
)

type Datasource struct {
	ID     string `bson:"uuid"`
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
		RawUrl: req.Name,
	}, nil
}

type FindAnagramsRequest struct {
	Datasource Datasource `bson:"datasource"`
	MaxWords   int        `bson:"maxWords"`
	MaxChars   int        `bson:"maxChars"`
}

func NewAnagramSettings(req *FindAnagramsRequest) {

}

type BitWeights struct {
	EncodedBits uint32  `bson:"encodedBits"`
	Weights     [26]int `bson:"weights"`
}

func NewBitWeights(s string) *BitWeights {
	bw := &BitWeights{}
	alphabet := "esiarntolcdugpmhbyfvkwzxjq"
	for i, char := range alphabet {
		if strings.ContainsRune(s, char) {
			bw.EncodedBits |= 1 << i
			bw.Weights[i] = strings.Count(s, string(char))
		}
	}

	return bw
}

func (bw *BitWeights) Combine(other *BitWeights) *BitWeights {
	newBW := &BitWeights{
		EncodedBits: bw.EncodedBits | other.EncodedBits,
	}

	for i := 0; i < 26; i++ {
		newBW.Weights[i] = bw.Weights[i] + other.Weights[i]
	}

	return newBW
}

func (bw *BitWeights) IsEquivalent(other *BitWeights) bool {
	return bw.EncodedBits == other.EncodedBits && bw.Weights == other.Weights
}

type AnagramEntry struct {
	ID           string     `bson:"uuid"`
	Encoded      BitWeights `bson:"encoded"`
	Anagrams     []string   `bson:"anagrams"`
	Combinations []int      `bson:"combinations"`
}

func NewAnagramEntry(s string) *AnagramEntry {
	bw := NewBitWeights(s)
	return &AnagramEntry{
		Encoded:      *bw,
		Anagrams:     []string{},
		Combinations: []int{1},
	}
}
