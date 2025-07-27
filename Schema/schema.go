package schema

import (
	"fmt"
	types "github.com/Abdullah05-js/goose/Types"
	"go.mongodb.org/mongo-driver/v2/bson"
	"reflect"
)

type Schema struct {
	Fields types.SchemaOptions
}

func NewSchema(fields types.SchemaOptions) (*Schema, error) {

	for key, opts := range fields {
		if opts.Default != nil && opts.Required {
			return nil, fmt.Errorf("can't use Default option with Required at the same time in field: %s", key)
		}
	}
	return &Schema{Fields: fields}, nil
}

func ValidateSchema[T types.ModelType](schema Schema, document T) (bson.M, error) {

	switch v := any(document).(type) {
	case bson.M:

		for key, opts := range schema.Fields {

			value, exists := v[key]

			if !exists {

				if opts.Required {
					return nil, fmt.Errorf("missing required field: %s", key)
				}

				if opts.Default != nil {
					v[key] = opts.Default
				}
				continue
			}

			if reflect.TypeOf(value).Kind() != opts.Type {
				return nil, fmt.Errorf("field %s expected type %s but got %s", key, opts.Type.String(), reflect.TypeOf(value).Kind().String())
			}

			if opts.Validate != nil {
				if err := opts.Validate(value); err != nil {
					return nil, fmt.Errorf("validation failed on field %s: %v", key, err)
				}
			}

		}

		return v, nil

	default:

		val := reflect.ValueOf(document)
		typ := reflect.TypeOf(document)

		newBson := make(bson.M)

		for i := 0; i < typ.NumField(); i++ {
			field := typ.Field(i)

			if !field.IsExported() { // if the filed is private skip
				continue
			}

			fieldName := field.Tag.Get("bson")

			if fieldName == "" {
				fieldName = field.Name
			}

			value := val.Field(i).Interface()
			opts, exists := schema.Fields[fieldName]

			if !exists {
				return nil, fmt.Errorf("missing Schema field: %s", field.Name)
			}

			if reflect.ValueOf(value).IsZero() {
				if opts.Required {
					return nil, fmt.Errorf("missing required field: %s", field.Name)
				}
				if opts.Default != nil {
					value = opts.Default
				}
			}

			if t := reflect.TypeOf(value); t != nil && t.Kind() != opts.Type {
				return nil, fmt.Errorf("field %s expected type %s but got %s", field.Name, opts.Type.String(), t.Kind().String())
			}

			if opts.Validate != nil {
				if err := opts.Validate(value); err != nil {
					return nil, fmt.Errorf("validation failed on field %s: %v", field.Name, err)
				}
			}
			newBson[fieldName] = value
		}
		return newBson, nil
	}
}
