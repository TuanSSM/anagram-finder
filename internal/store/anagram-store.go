package store

import (
	"context"
	"fmt"

	"github.com/tuanssm/anagram-finder/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnagramStorer interface {
	Insert(context.Context, string, *types.AnagramEntry) error
	BulkInsert(context.Context, string, []*types.AnagramEntry) error
	GetByBitWeights(context.Context, string) (*types.AnagramEntry, error)
	GetAll(context.Context) ([]*types.AnagramEntry, error)
	Append(context.Context, string, string, int, []string) error
	GetArrayAtIndex(context.Context, string, string, int) ([]string, error)
}

type AnagramStore struct {
	db *mongo.Database
}

func NewAnagramStore(db *mongo.Database) *AnagramStore {
	return &AnagramStore{
		db: db,
	}
}

func (s *AnagramStore) Insert(ctx context.Context, coll string, a *types.AnagramEntry) error {
	filter := bson.M{"encoded": a.Encoded}

	update := bson.M{
		"$setOnInsert": bson.M{
			"encoded":      a.Encoded,
			"combinations": a.Combinations,
		},
		"$addToSet": bson.M{
			"anagrams": bson.M{
				"$each": a.Anagrams,
			},
		},
	}

	upsert := true
	opts := options.Update().SetUpsert(upsert)

	res, err := s.db.Collection(coll).UpdateOne(ctx, filter, update, opts)
	if err != nil {
		return err
	}

	if res.UpsertedID != nil {
		a.ID = res.UpsertedID.(primitive.ObjectID).Hex()
	}

	return nil
}

func (s *AnagramStore) BulkInsert(ctx context.Context, coll string, anagrams []*types.AnagramEntry) error {
	documents := make([]interface{}, len(anagrams))
	for i, a := range anagrams {
		documents[i] = bson.M{
			"encoded":      a.Encoded,
			"anagrams":     a.Anagrams,
			"combinations": a.Combinations,
		}
	}

	_, err := s.db.Collection(coll).InsertMany(ctx, documents)
	return err
}

func (s *AnagramStore) GetAll(ctx context.Context, coll string) ([]*types.AnagramEntry, error) {
	cursor, err := s.db.Collection(coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	anagrams := []*types.AnagramEntry{}
	err = cursor.All(ctx, &anagrams)
	return anagrams, err
}

func (s *AnagramStore) Append(ctx context.Context, coll string, id string, index int, value []string) error {
	filter := bson.M{"uuid": id}

	curAnagram, err := s.GetArrayAtIndex(ctx, coll, id, index-1)
	if err != nil {
		return err
	}

	var update bson.M
	if curAnagram != nil {
		update = bson.M{
			"$push": bson.M{
				fmt.Sprintf("rawUrl.%d", index): bson.M{
					"$each": value,
				},
			},
		}
	} else {
		// Fill with empty arrays until the index
		emptyArrs := make([][]string, index+1)
		emptyArrs[index] = value
		update = bson.M{
			"$set": bson.M{
				"rawUrl": emptyArrs,
			},
		}
	}

	// Apply the update operation to the database
	_, err = s.db.Collection(coll).UpdateOne(ctx, filter, update)
	return err
}

func (s *AnagramStore) GetArrayAtIndex(ctx context.Context, coll string, id string, index int) ([]string, error) {
	// Create a filter to find the document
	filter := bson.M{"uuid": id}

	// Use projection with the $slice operator to get the array at the given index
	projection := bson.M{
		"anagrams": bson.M{
			"$slice": []int{index, 1},
		},
	}

	var result struct {
		Anagrams [][]string `bson:"anagrams"`
	}
	err := s.db.Collection(coll).FindOne(ctx, filter, options.FindOne().SetProjection(projection)).Decode(&result)
	if err != nil {
		return nil, err
	}

	if len(result.Anagrams) > 0 {
		return result.Anagrams[0], nil
	}

	// If there is no array at the given index, return nil
	return nil, nil
}
