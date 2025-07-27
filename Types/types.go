package types

import (
	"reflect"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type OnlyStruct interface {
	struct{}
}

type FieldOptions struct {
	Required bool
	Default  interface{}
	Type     reflect.Kind
	Validate func(interface{}) error // optional custom validator
}

type SchemaOptions map[string]FieldOptions

type MongoStore struct {
	Client *mongo.Client
	DB     *mongo.Database
}

type ModelType interface {
	OnlyStruct | bson.M
}
