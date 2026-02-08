package database

import (
	"context"
	"fly-go/internal/utils"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func (m *MongoDB) Find(ctx context.Context, collectionName string, query utils.BaseQuery) Rows {
	filter := map[string]interface{}{}
	if query.Search != "" {
		filter["name"] = map[string]interface{}{"$regex": query.Search}
	}
	opts := options.Find()
	if query.OrderBy != "" {
		opts.SetSort(map[string]int{query.OrderBy: 1})
	}
	if query.Page > 0 && query.Size > 0 {
		opts.SetSkip(int64((query.Page - 1) * query.Size))
		opts.SetLimit(int64(query.Size))
	}
	cur, err := m.Collection(collectionName).Find(ctx, filter, opts)
	if err != nil {
		return nil
	}
	results := Rows{}
	err = cur.All(ctx, &results)
	if err != nil {
		return nil
	}
	return results
}

func (m *MongoDB) FindOne(ctx context.Context, collectionName string, filter interface{}, opts ...*options.FindOneOptions) *mongo.SingleResult {
	return m.Collection(collectionName).FindOne(ctx, filter, opts...)
}
