package main

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnagramStorer interface {
	Insert(context.Context, string, *AnagramEntry) error
	GetByBitWeights(context.Context, string) (*AnagramEntry, error)
	GetAll(context.Context) ([]*AnagramEntry, error)
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

func (s *AnagramStore) Insert(ctx context.Context, coll string, a *AnagramEntry) error {
	res, err := s.db.Collection(coll).InsertOne(ctx, a)
	if err != nil {
		return err
	}
	a.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return err
}

func (s *AnagramStore) GetAll(ctx context.Context, coll string) ([]*AnagramEntry, error) {
	cursor, err := s.db.Collection(coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	anagrams := []*AnagramEntry{}
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
