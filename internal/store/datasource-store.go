package store

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"

	"github.com/tuanssm/anagram-finder/internal/types"
)

type DatasourceStorer interface {
	Insert(context.Context, *types.Datasource) (*types.Datasource, error)
	GetByID(context.Context, string) (*types.Datasource, error)
	GetAll(context.Context) ([]*types.Datasource, error)
}

type DatasourceStore struct {
	db   *mongo.Database
	coll string
}

func NewDatasourceStore(db *mongo.Database) *DatasourceStore {
	return &DatasourceStore{
		db:   db,
		coll: "datasources",
	}
}

func (s *DatasourceStore) Insert(ctx context.Context, d *types.Datasource) (*types.Datasource, error) {
	d.ID = primitive.NewObjectID().Hex()
	_, err := s.db.Collection(s.coll).InsertOne(ctx, d)
	if err != nil {
		return nil, err
	}

	return d, err
}

func (s *DatasourceStore) GetAll(ctx context.Context) ([]*types.Datasource, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	datasources := []*types.Datasource{}
	err = cursor.All(ctx, &datasources)
	return datasources, err
}

func (s *DatasourceStore) GetByID(ctx context.Context, id string) (*types.Datasource, error) {
	var (
		objID, _ = primitive.ObjectIDFromHex(id)
		res      = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		d        = &types.Datasource{}
		err      = res.Decode(d)
	)
	return d, err
}
