# Goose 🦆

A lightweight, type-safe MongoDB Object Document Mapper (ODM) for Go with schema validation, flexible data modeling, and relationship population.

## Features

- 🔒 **Type Safety**: Full generic support for type-safe operations
- 📋 **Schema Validation**: Built-in field validation with custom validators
- 🔗 **Relationships**: Support for document references with population
- 🏗️ **Flexible Models**: Support for both structs and `bson.M` documents
- 🚀 **Easy Setup**: Simple connection management
- ⚡ **Performance**: Lightweight wrapper around the official MongoDB Go driver
- 🛡️ **Error Handling**: Comprehensive error messages for debugging
- 🔍 **Query Chaining**: Fluent API for query operations with populate support

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
    "profile": {
        Required: false,
        Type:     reflect.String, // ObjectID as string
        Ref:      "profiles",     // Reference to profiles collection
    },
})
```

### 3. Create and Use Models

#### Using Structs

```go
type User struct {
    Name    string             `bson:"name"`
    Email   string             `bson:"email"`
    Age     int                `bson:"age"`
    Profile bson.ObjectID      `bson:"profile,omitempty"`
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

result, err := userModel.InsertOne(ctx, user)
if err != nil {
    log.Fatal(err)
}

// Find documents with query chaining
query := userModel.Find(ctx, bson.M{"age": bson.M{"$gte": 18}})
users, _, err := query.Result()
if err != nil {
    log.Fatal(err)
}

// Find one document with population
query = userModel.FindOne(ctx, bson.M{"email": "john@example.com"}).Populate("profile")
user, _, err := query.Result()
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

result, err := userModel.InsertOne(ctx, userData)
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
    Ref      string                         // Collection reference for population
    Validate func(interface{}) error        // Custom validation function
}
```

### Example Schema with All Options

```go
postSchema, err := schema.NewSchema(types.SchemaOptions{
    "title": {
        Required: true,
        Type:     reflect.String,
        Validate: func(v interface{}) error {
            title := v.(string)
            if len(title) < 5 {
                return fmt.Errorf("title must be at least 5 characters")
            }
            return nil
        },
    },
    "content": {
        Required: true,
        Type:     reflect.String,
    },
    "author": {
        Required: true,
        Type:     reflect.String, // ObjectID as string
        Ref:      "users",        // Reference to users collection
    },
    "category": {
        Required: false,
        Type:     reflect.String,
        Default:  "general",
    },
    "published": {
        Required: false,
        Type:     reflect.Bool,
        Default:  false,
    },
})
```

## Query System

The new query system provides a fluent API for chaining operations:

### Basic Queries

```go
// Find many documents
query := userModel.Find(ctx, bson.M{"age": bson.M{"$gte": 18}})
_, users, err := query.Result()
if err != nil {
    log.Fatal(err)
}

// Find one document
query = userModel.FindOne(ctx, bson.M{"email": "john@example.com"})
user, _, err := query.Result()
if err != nil {
    log.Fatal(err)
}
```

### Population (Relationships)

Population allows you to automatically fetch referenced documents:

```go
// Populate a single reference
query := userModel.FindOne(ctx, bson.M{"_id": userID}).Populate("profile")
user, _, err := query.Result()
if err != nil {
    log.Fatal(err)
}

// Populate multiple documents
query = postModel.Find(ctx, bson.M{}).Populate("author")
_, posts, err := query.Result()
if err != nil {
    log.Fatal(err)
}
// Each post will have the full author document instead of just the ObjectID
```

### Query Result Handling

The `Result()` method returns different values based on the query type:

```go
// For FindOne queries
singleResult, _, err := query.Result()

// For Find queries  
_, manyResults, err := query.Result()
```

## API Reference

### Connection Management

#### `Connect[T *options.ClientOptions | string](ClientOptions T, dataBaseName string) error`
Establishes connection to MongoDB database.

#### `DisConnect(ctx context.Context) error`
Closes the MongoDB connection.

### Model Operations

#### `InsertOne(ctx context.Context, data T, opts ...options.Lister[options.InsertOneOptions]) (*mongo.InsertOneResult, error)`
Validates and inserts a single document.

#### `FindOne(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOneOptions]) *Query[T]`
Finds a single document matching the query and returns a Query object for chaining.

#### `Find(ctx context.Context, query bson.M, opts ...options.Lister[options.FindOptions]) *Query[T]`
Finds multiple documents matching the query and returns a Query object for chaining.

### Query Operations

#### `Populate(fieldName string) *Query[T]`
Populates a referenced field with the actual document from the referenced collection.

#### `Result() (T, []T, error)`
Executes the query and returns the results. For single queries, returns `(result, nil, error)`. For multiple queries, returns `(zeroValue, results, error)`.

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
- **Population errors**: `"populate: field fieldName is not an ObjectID"`
- **Reference errors**: `"Ref in field Options not defined"`

## Complete Example with Relationships

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
    ID    bson.ObjectID `bson:"_id,omitempty"`
    Name  string        `bson:"name"`
    Email string        `bson:"email"`
}

type Post struct {
    ID      bson.ObjectID `bson:"_id,omitempty"`
    Title   string        `bson:"title"`
    Content string        `bson:"content"`
    Author  bson.ObjectID `bson:"author"`
}

func main() {
    ctx := context.Background()
    
    // Connect to MongoDB
    err := goose.Connect("mongodb://localhost:27017", "blogDB")
    if err != nil {
        log.Fatal("Connection failed:", err)
    }
    defer goose.DisConnect(ctx)
    
    // Create user schema
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
    })
    if err != nil {
        log.Fatal("User schema creation failed:", err)
    }
    
    // Create post schema with user reference
    postSchema, err := schema.NewSchema(types.SchemaOptions{
        "title": {
            Required: true,
            Type:     reflect.String,
        },
        "content": {
            Required: true,
            Type:     reflect.String,
        },
        "author": {
            Required: true,
            Type:     reflect.String, // ObjectID field
            Ref:      "users",        // Reference to users collection
        },
    })
    if err != nil {
        log.Fatal("Post schema creation failed:", err)
    }
    
    // Create models
    userModel, err := model.NewModel[User](ctx, "users", userSchema)
    if err != nil {
        log.Fatal("User model creation failed:", err)
    }
    
    postModel, err := model.NewModel[Post](ctx, "posts", postSchema)
    if err != nil {
        log.Fatal("Post model creation failed:", err)
    }
    
    // Insert a user
    user := User{
        Name:  "Alice Johnson",
        Email: "alice@example.com",
    }
    
    userResult, err := userModel.InsertOne(ctx, user)
    if err != nil {
        log.Fatal("User insert failed:", err)
    }
    
    userID := userResult.InsertedID.(bson.ObjectID)
    fmt.Printf("Inserted user with ID: %v\n", userID)
    
    // Insert a post with user reference
    post := Post{
        Title:   "My First Blog Post",
        Content: "This is the content of my first blog post.",
        Author:  userID,
    }
    
    postResult, err := postModel.InsertOne(ctx, post)
    if err != nil {
        log.Fatal("Post insert failed:", err)
    }
    
    fmt.Printf("Inserted post with ID: %v\n", postResult.InsertedID)
    
    // Find posts with populated author information
    query := postModel.Find(ctx, bson.M{}).Populate("author")
    _, posts, err := query.Result()
    if err != nil {
        log.Fatal("Find with populate failed:", err)
    }
    
    fmt.Printf("Found %d posts\n", len(posts))
    for _, p := range posts {
        fmt.Printf("Post: %+v\n", p)
        // The author field now contains the full user document instead of just ObjectID
    }
}
```

## Best Practices

### Query Optimization
When using `Find()`, consider setting a limit to avoid large result sets:

```go
opts := options.Find().SetLimit(10)
query := userModel.Find(ctx, bson.M{}, opts)
```

### Population Performance
Population executes additional database queries. Use it judiciously:

```go
// Good: Only populate when you need the referenced data
query := postModel.FindOne(ctx, bson.M{"_id": postID}).Populate("author")

// Consider: Do you really need to populate all posts?
query = postModel.Find(ctx, bson.M{}).Populate("author")
```

### Error Handling
Always handle errors from the Result() method:

```go
query := userModel.FindOne(ctx, bson.M{"email": "nonexistent@example.com"})
user, _, err := query.Result()
if err != nil {
    // Handle specific errors
    if strings.Contains(err.Error(), "no document found") {
        // User not found
        return nil
    }
    // Other error
    return err
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