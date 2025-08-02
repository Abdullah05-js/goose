package utils

import "go.mongodb.org/mongo-driver/v2/bson"

func ToT[T any](m bson.M) (T, error) {
	var result T
	data, err := bson.Marshal(m)
	if err != nil {
		return result, err
	}
	err = bson.Unmarshal(data, &result)
	return result, err
}
