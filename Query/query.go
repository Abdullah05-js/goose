package query

import (
	"context"
	"fmt"

	"github.com/Abdullah05-js/goose"
	shared "github.com/Abdullah05-js/goose/Shared"
	types "github.com/Abdullah05-js/goose/Types"
	"github.com/Abdullah05-js/goose/utils"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type ResultType string

const (
	ResultMany   ResultType = "many"
	ResultSingle ResultType = "single"
)

type Query[T types.ModelType] struct {
	Model        *shared.BaseModel[T]
	Type         ResultType
	ManyResult   []bson.M
	SingleResult bson.M
	Err          error
	Ctx          context.Context
}

func NewQuery[T types.ModelType](model *shared.BaseModel[T], queryType ResultType, ctx context.Context) *Query[T] {
	return &Query[T]{
		Model: model,
		Type:  queryType,
		Ctx:   ctx,
	}
}

func (q *Query[T]) Populate(filedName string) *Query[T] {
	if q.Err != nil {
		return q
	}
	// Schemadan ref all
	filedOpts, ok := q.Model.Schema.Fields[filedName]
	if !ok {
		q.Err = fmt.Errorf("from populate func: filedName isnt valid")
		return q
	}
	targetRef := filedOpts.Ref
	var ZeroString string
	if targetRef == ZeroString {
		q.Err = fmt.Errorf("from populate func: Ref in field Options not defined")
		return q
	}

	targetCollection := goose.Base.DB.Collection(targetRef)

	switch q.Type {
	case ResultSingle:
		id, ok := q.SingleResult[filedName].(bson.ObjectID)
		if !ok {
			q.Err = fmt.Errorf("populate: field %q is not an ObjectID", filedName)
			return q
		}
		var targetResult bson.M
		if err := targetCollection.FindOne(q.Ctx, bson.M{"_id": id}).Decode(&targetResult); err != nil {
			q.Err = fmt.Errorf("populate: failed to find document with _id %v: %w", id, err)
			return q
		}
		q.SingleResult[filedName] = targetResult

	case ResultMany:
		for i, obj := range q.ManyResult {
			id, ok := obj[filedName].(bson.ObjectID)
			if !ok {
				q.Err = fmt.Errorf("populate: field %q in document #%d is not an ObjectID", filedName, i)
				return q
			}
			var targetResult bson.M
			if err := targetCollection.FindOne(q.Ctx, bson.M{"_id": id}).Decode(&targetResult); err != nil {
				q.Err = fmt.Errorf("populate: failed to find document with _id %v in document #%d: %w", id, i, err)
				return q
			}
			obj[filedName] = targetResult
		}

	default:
		q.Err = fmt.Errorf("populate: unknown ResultType %q", q.Type)
		return q
	}
	return q
}

func (q *Query[T]) Result() (T, []T, error) {
	var zeroValue T

	if q.Err != nil {
		return zeroValue, nil, q.Err
	}

	switch q.Type {
	case ResultSingle:
		result, err := utils.ToT[T](q.SingleResult)
		if err != nil {
			return zeroValue, nil, err
		}
		return result, nil, nil

	case ResultMany:
		resultArr := make([]T, len(q.ManyResult))
		for i, v := range q.ManyResult {
			result, err := utils.ToT[T](v)
			if err != nil {
				return zeroValue, nil, err
			}
			resultArr[i] = result
		}
		return zeroValue, resultArr, nil

	default:
		err := fmt.Errorf("result: unknown ResultType %q", q.Type)
		return zeroValue, nil, err
	}
}
