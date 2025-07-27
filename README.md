# Goose 🦆

A lightweight, type-safe MongoDB Object Document Mapper (ODM) for Go with schema validation and flexible data modeling.

## Features

- 🔒 **Type Safety**: Full generic support for type-safe operations
- 📋 **Schema Validation**: Built-in field validation with custom validators
- 🏗️ **Flexible Models**: Support for both structs and `bson.M` documents
- 🚀 **Easy Setup**: Simple connection management
- ⚡ **Performance**: Lightweight wrapper around the official MongoDB Go driver
- 🛡️ **Error Handling**: Comprehensive error messages for debugging

## Installation

```bash
go get github.com/Abdullah05-js/goose
```

## Quick Start

### 1. Connect to MongoDB

```go
package main

import (
    "context"
    "log"
    "github.com/Abdullah05-js/goose"
    "go.mongodb.org/mongo-driver/v2/mongo/options"
)

func main() {
    // Using connection string
    err := goose.Connect("mongodb://localhost:27017", "myDatabase")
    if err != nil {
        log.Fatal(err)
    }
    
    // Or using ClientOptions
    clientOpts := options.Client().ApplyURI("mongodb://localhost:27017")
    err = goose.Connect(clientOpts, "myDatabase")
    if err != nil {
        log.Fatal(err)
    }
    
    defer goose.DisConnect(context.Background())
}
```

### 2. Define a Schema

```go
import (
    "reflect"
    "github.com/Abdullah05-js/goose/Schema"
    "github.com/Abdullah05-js/goose/Types"
)

// Create schema with validation rules
userSchema, err := schema.NewSchema(types.SchemaOptions{
    "name": {
        Required: true,
        Type:     reflect.String,
    },
    "email": {
        Required: true,
        Type:     reflect.String,
        Validate: func(v interface{}) error {
            email := v.(string)
            if !strings.Contains(email, "@") {
                return fmt.Errorf("invalid email format")
            }
            return nil
        },
    },
    "age": {
        Required: false,
        Type:     reflect.Int,
        Default:  18,
    },
})
```

### 3. Create and Use Models

#### Using Structs

```go
type User struct {
    Name  string `bson:"name"`
    Email string `bson:"email"`
    Age   int    `bson:"age"`
}

// Create model
userModel, err := model.NewModel[User](ctx, "users", userSchema)
if err != nil {
    log.Fatal(err)
}

// Insert a document
user := User{
    Name:  "John Doe",
    Email: "john@example.com",
    Age:   25,
}

result, err := userModel.InserOne(ctx, user)
if err != nil {
    log.Fatal(err)
}

// Find documents
users, err := userModel.Find(ctx, bson.M{"age": bson.M{"$gte": 18}})
if err != nil {
    log.Fatal(err)
}

// Find one document
user, err := userModel.FindOne(ctx, bson.M{"email": "john@example.com"})
if err != nil {
    log.Fatal(err)
}
```

#### Using bson.M

```go
// Create model with bson.M
userModel, err := model.NewModel[bson.M](ctx, "users", userSchema)
if err != nil {
    log.Fatal(err)
}

// Insert document
userData := bson.M{
    "name":  "Jane Doe",
    "email": "jane@example.com",
    "age":   30,
}

result, err := userModel.InserOne(ctx, userData)
if err != nil {
    log.Fatal(err)
}
```

## Schema Options

Configure field validation and behavior with `FieldOptions`:

```go
type FieldOptions struct {
    Required bool                           // Field is required
    Default  interface{}                    // Default value if not provided
    Type     reflect.Kind                   // Expected data type
    Validate func(interface{}) error        // Custom validation function
}
```

### Example Schema with All Options

```go
productSchema, err := schema.NewSchema(types.SchemaOptions{
    "name": {
        Required: true,
        Type:     reflect.String,
        Validate: func(v interface{}) error {
            name := v.(string)
            if len(name) < 3 {
                return fmt.Errorf("name must be at least 3 characters")
            }
            return nil
        },
    },
    "price": {
        Required: true,
        Type:     reflect.Float64,
        Validate: func(v interface{}) error {
            price := v.(float64)
            if price <= 0 {
                return fmt.Errorf("price must be positive")
            }
            return nil
        },
    },
    "category": {
        Required: false,
        Type:     reflect.String,
        Default:  "uncategorized",
    },
    "inStock": {
        Required: false,
        Type:     reflect.Bool,
        Default:  true,
    },
})
```

## API Reference

### Connection Management

#### `Connect[T *options.ClientOptions | string](ClientOptions T, dataBaseName string) error`
Establishes connection to MongoDB database.

#### `DisConnect(ctx context.Context) error`
Closes the MongoDB connection.

### Model Operations

#### `InserOne(ctx context.Context, data T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)`
Validates and inserts a single document.

#### `FindOne(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOneOptions]) (T, error)`
Finds and returns a single document matching the query.

#### `Find(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOptions]) ([]T, error)`
Finds and returns multiple documents matching the query.

### Schema Validation

#### `NewSchema(fields types.SchemaOptions) (*Schema, error)`
Creates a new schema with validation rules.

#### `ValidateSchema[T types.ModelType](schema Schema, document T) (bson.M, error)`
Validates a document against the schema.

## Error Handling

Goose provides detailed error messages for common scenarios:

- **Missing required fields**: `"missing required field: fieldName"`
- **Type mismatches**: `"field fieldName expected type string but got int"`
- **Validation failures**: `"validation failed on field fieldName: custom error"`
- **Schema conflicts**: `"can't use Default option with Required at the same time"`

## Examples

### Complete Example

```go
package main

import (
    "context"
    "fmt"
    "log"
    "reflect"
    "strings"
    
    "github.com/Abdullah05-js/goose"
    "github.com/Abdullah05-js/goose/Schema"
    "github.com/Abdullah05-js/goose/Types"
    "github.com/Abdullah05-js/goose/model"
    "go.mongodb.org/mongo-driver/v2/bson"
)

type User struct {
    Name  string `bson:"name"`
    Email string `bson:"email"`
    Age   int    `bson:"age"`
}

func main() {
    ctx := context.Background()
    
    // Connect to MongoDB
    err := goose.Connect("mongodb://localhost:27017", "testDB")
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer goose.DisConnect(ctx)
    
    // Create schema
    userSchema, err := schema.NewSchema(types.SchemaOptions{
        "name": {
            Required: true,
            Type:     reflect.String,
        },
        "email": {
            Required: true,
            Type:     reflect.String,
            Validate: func(v interface{}) error {
                email := v.(string)
                if !strings.Contains(email, "@") {
                    return fmt.Errorf("invalid email format")
                }
                return nil
            },
        },
        "age": {
            Required: false,
            Type:     reflect.Int,
            Default:  18,
        },
    })
    if err != nil {
        log.Fatal("Schema creation failed:", err)
    }
    
    // Create model
    userModel, err := model.NewModel[User](ctx, "users", userSchema)
    if err != nil {
        log.Fatal("Model creation failed:", err)
    }
    
    // Insert user
    user := User{
        Name:  "Alice Johnson",
        Email: "alice@example.com",
        Age:   28,
    }
    
    result, err := userModel.InserOne(ctx, user)
    if err != nil {
        log.Fatal("Insert failed:", err)
    }
    
    fmt.Printf("Inserted document with ID: %v\n", result.InsertedID)
    
    // Find users
    users, err := userModel.Find(ctx, bson.M{})
    if err != nil {
        log.Fatal("Find failed:", err)
    }
    
    fmt.Printf("Found %d users\n", len(users))
    for _, u := range users {
        fmt.Printf("User: %+v\n", u)
    }
}
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Support

If you have any questions or need help, please:
- Open an issue on GitHub
- Check the [documentation](https://github.com/Abdullah05-js/goose)
- Review the examples in the `/examples` directory

---

**Goose** - Making MongoDB operations in Go as smooth as a goose gliding on water! 🦆