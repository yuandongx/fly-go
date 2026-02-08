package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) DeleteOne(ctx context.Context, collectionName string, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	collection, err := m.Collection(collectionName)
	if err != nil {
		return 0, err
	}
	res, err := collection.DeleteOne(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}

func (m *MongoDB) DeleteMany(ctx context.Context, collectionName string, filter interface{}, opts ...*options.DeleteOptions) (int64, error) {
	collection, err := m.Collection(collectionName)
	if err != nil {
		return 0, err
	}
	res, err := collection.DeleteMany(ctx, filter, opts...)
	if err != nil {
		return 0, err
	}
	return res.DeletedCount, nil
}
