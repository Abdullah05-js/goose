package utils

import (
	"context"
	"fmt"
	"github.com/Abdullah05-js/goose"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

func GetCollection(ctx context.Context, collectionName string) (*mongo.Collection, error) {
	collections, err := goose.Base.DB.ListCollectionNames(ctx, bson.M{"name": collectionName})
	if err != nil {
		return nil, fmt.Errorf("failed to list collections:%w", err)
	}

	if len(collections) == 0 {
		err := goose.Base.DB.CreateCollection(ctx, collectionName)
		if err != nil {
			return nil, fmt.Errorf("failed to create collection:%w", err)
		}
	}
	return goose.Base.DB.Collection(collectionName), nil
}