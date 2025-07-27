package goose

import (
	"context"
	"fmt"
	types "github.com/Abdullah05-js/goose/Types"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

var Base *types.MongoStore

func Connect[T *options.ClientOptions | string](ClientOptions T, dataBaseName string) error {
	Base = new(types.MongoStore)
	switch v := any(ClientOptions).(type) {
	case string:
		client, err := mongo.Connect(options.Client().ApplyURI(v))
		if err != nil {
			return err
		}
		Base.Client = client
	case *options.ClientOptions:
		client, err := mongo.Connect(v)
		if err != nil {
			return err
		}
		Base.Client = client
	default:
		return fmt.Errorf("the type can only be options.ClientOptions or string")
	}
	Base.DB = Base.Client.Database(dataBaseName)
	return nil
}

func DisConnect(ctx context.Context) error {
	return Base.Client.Disconnect(ctx)
}
