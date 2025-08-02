package types

import (
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"reflect"
)

type OnlyStruct interface {
	struct{}
}

type FieldOptions struct {
	Required bool
	Default  interface{}
	Type     reflect.Kind
	Ref      string
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
