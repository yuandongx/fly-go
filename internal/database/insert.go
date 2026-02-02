package database

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) InsertOne(ctx context.Context, collectionName string, document interface{}, opts ...*options.InsertOneOptions) (interface{}, error) {
	res, err := m.Collection(collectionName).InsertOne(ctx, document, opts...)
	if err != nil {
		return nil, err
	}
	return res.InsertedID, nil
}

func (m *MongoDB) InsertMany(ctx context.Context, collectionName string, documents interface{}, opts ...*options.InsertManyOptions) ([]interface{}, error) {
	res, err := m.Collection(collectionName).InsertMany(ctx, documents.([]interface{}), opts...)
	if err != nil {
		return nil, err
	}
	return res.InsertedIDs, nil
}
