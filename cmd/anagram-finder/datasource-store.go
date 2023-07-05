package main

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type DatasourceStorer interface {
	Insert(context.Context, *Datasource) error
	GetByID(context.Context, string) (*Datasource, error)
	GetAll(context.Context) ([]*Datasource, error)
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

func (s *DatasourceStore) Insert(ctx context.Context, d *Datasource) error {
	res, err := s.db.Collection(s.coll).InsertOne(ctx, d)
	if err != nil {
		return err
	}
	d.ID = res.InsertedID.(primitive.ObjectID).Hex()

	return err
}

func (s *DatasourceStore) GetAll(ctx context.Context) ([]*Datasource, error) {
	cursor, err := s.db.Collection(s.coll).Find(ctx, map[string]any{})
	if err != nil {
		return nil, err
	}

	datasources := []*Datasource{}
	err = cursor.All(ctx, &datasources)
	return datasources, err
}

func (s *DatasourceStore) GetByID(ctx context.Context, id string) (*Datasource, error) {
	var (
		objID, _ = primitive.ObjectIDFromHex(id)
		res      = s.db.Collection(s.coll).FindOne(ctx, bson.M{"_id": objID})
		d        = &Datasource{}
		err      = res.Decode(d)
	)
	return d, err
}
