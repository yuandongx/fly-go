package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) Find(ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOptions) (*mongo.Cursor, error) {
	return m.Collection(collectionName).Find(ctx, filter, opts...)
}

func (m *MongoDB) FindOne(ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.Collection(collectionName).FindOne(ctx, filter, opts...)
}
