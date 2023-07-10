package store

import (
	"context"

	"github.com/tuanssm/anagram-finder/internal/types"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type AnagramStorer interface {
	Insert(context.Context, string, *types.AnagramEntry) error
	BulkInsert(context.Context, string, []*types.AnagramEntry) error
	BulkUpsert(context.Context, string, []*types.AnagramEntry) error
	GetAll(context.Context) ([]*types.AnagramEntry, error)
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
			"encoded": a.Encoded,
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
			"encoded":  a.Encoded,
			"anagrams": a.Anagrams,
		}
	}

	_, err := s.db.Collection(coll).InsertMany(ctx, documents)
	return err
}

func (s *AnagramStore) BulkUpsert(ctx context.Context, coll string, anagrams []*types.AnagramEntry) error {
	var models []mongo.WriteModel

	for _, a := range anagrams {
		filter := bson.M{"encoded": a.Encoded}

		update := bson.M{
			"$setOnInsert": bson.M{
				"encoded": a.Encoded,
			},
			"$addToSet": bson.M{
				"anagrams": bson.M{
					"$each": a.Anagrams,
				},
			},
		}

		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}

	_, err := s.db.Collection(coll).BulkWrite(ctx, models)
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
