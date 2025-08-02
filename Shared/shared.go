package shared

import (
	schema "github.com/Abdullah05-js/goose/Schema"
	types "github.com/Abdullah05-js/goose/Types"
	"go.mongodb.org/mongo-driver/v2/mongo"
)

type BaseModel[T types.ModelType] struct {
	Collection *mongo.Collection
	Schema     *schema.Schema
}


