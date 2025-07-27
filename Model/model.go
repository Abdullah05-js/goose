package model

import (
	"context"
	"fmt"
	s "github.com/Abdullah05-js/goose/Schema"
	types "github.com/Abdullah05-js/goose/Types"
	"github.com/Abdullah05-js/goose/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Model[T types.ModelType] struct {
	Collection *mongo.Collection
	Schema     *s.Schema
}

func newModel[T types.ModelType](ctx context.Context, collectionName string, schema *s.Schema) (*Model[T], error) {
	collection, err := utils.GetCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	return &Model[T]{Collection: collection, Schema: schema}, nil
}

func (model *Model[T]) InserOne(ctx context.Context, data T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
	val, err := s.ValidateSchema(*model.Schema, data)
	if err != nil {
		return nil, err
	}
	result, errInsertOne := model.Collection.InsertOne(ctx, val, opts...)
	if errInsertOne != nil {
		return nil, errInsertOne
	}
	return result, nil
}

func (model *Model[T]) FindOne(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOneOptions]) (T, error) {
	var result T
	err := model.Collection.FindOne(ctx, query, opts...).Decode(&result)

	if err != nil {
		var Zero T
		if err == mongo.ErrNoDocuments {
			return Zero, fmt.Errorf("no document found by the query: %s", err)
		}
		return Zero, err
	}
	return result, nil
}

func (model *Model[T]) Find(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error) {
	var result []T
	cursor, err := model.Collection.Find(ctx, query, opts...)
	if err != nil {
		return nil, err
	}
	defer cursor.Close(ctx)
	if err := cursor.All(ctx, &result); err != nil {
		return nil, err
	}
	return result, nil
}
