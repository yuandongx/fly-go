package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) UpdateOne(ctx context.Context, collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	collection, err := m.Collection(collectionName)
	if err != nil {
		return 0, err
	}
	res, err := collection.UpdateOne(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}

func (m *MongoDB) UpdateMany(ctx context.Context, collectionName string, filter interface{}, update interface{}, opts ...*options.UpdateOptions) (int64, error) {
	collection, err := m.Collection(collectionName)
	if err != nil {
		return 0, err
	}
	res, err := collection.UpdateMany(ctx, filter, update, opts...)
	if err != nil {
		return 0, err
	}
	return res.ModifiedCount, nil
}
