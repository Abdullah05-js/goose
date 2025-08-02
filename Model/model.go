package model

import (
	"context"
	"fmt"
	q "github.com/Abdullah05-js/goose/Query"
	s "github.com/Abdullah05-js/goose/Schema"
	shared "github.com/Abdullah05-js/goose/Shared"
	types "github.com/Abdullah05-js/goose/Types"
	"github.com/Abdullah05-js/goose/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type Model[T types.ModelType] struct {
	shared.BaseModel[T]
}

func newModel[T types.ModelType](ctx context.Context, collectionName string, schema *s.Schema) (*Model[T], error) {
	collection, err := utils.GetCollection(ctx, collectionName)
	if err != nil {
		return nil, err
	}
	return &Model[T]{shared.BaseModel[T]{Collection: collection, Schema: schema}}, nil
}

func (model *Model[T]) InsertOne(ctx context.Context, data T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error) {
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

func (model *Model[T]) FindOne(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOneOptions]) *q.Query[T] {
	qry := q.NewQuery(&model.BaseModel, q.ResultSingle, ctx)
	var result bson.M
	err := model.Collection.FindOne(ctx, query, opts...).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			qry.Err = fmt.Errorf("no document found by the query: %s", err)
			return qry
		}
		qry.Err = err
		return qry
	}
	qry.SingleResult = result
	return qry
}

// When using this function, consider setting a limit with options.SetLimit() to avoid large result sets.
func (model *Model[T]) Find(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOptions]) *q.Query[T] {
	qry := q.NewQuery(&model.BaseModel, q.ResultMany, ctx)

	cursor, err := model.Collection.Find(ctx, query, opts...)
	if err != nil {
		qry.Err = err
		return qry
	}
	defer cursor.Close(ctx)

	var result []bson.M
	if err := cursor.All(ctx, &result); err != nil {
		qry.Err = err
		return qry
	}

	qry.ManyResult = result
	return qry
}
