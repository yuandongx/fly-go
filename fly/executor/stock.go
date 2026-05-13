package executor

import (
	"context"
	"fmt"

	"fly-go/database"
	"fly-go/fly"
	"fly-go/fly/spider"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// StockExecutor 股票数据抓取执行器
type StockExecutor struct {
	db *database.MongoDB
}

func NewStockExecutor(db *database.MongoDB) *StockExecutor {
	return &StockExecutor{db: db}
}

func (e *StockExecutor) Name() string {
	return "stock"
}

func (e *StockExecutor) Execute(ctx context.Context, task *fly.Task) (fly.TaskResult, error) {
	result := fly.NewTaskResult()

	fmt.Printf("Fetching stock data...\n")

	// 获取股票数据
	data, err := spider.GetStockInfo()
	if err != nil {
		result.Status = fly.StatusError
		result.Message = fmt.Sprintf("Failed to fetch stock data: %v", err)
		return result, err
	}

	result.Data = data
	result.Message = fmt.Sprintf("Fetched %d stock records", len(data))

	// 保存到数据库
	if task.OutputColl != "" && len(data) > 0 {
		if err := e.saveToDatabase(ctx, task.OutputColl, data); err != nil {
			result.Message = fmt.Sprintf("Fetched %d records but failed to save: %v", len(data), err)
			return result, err
		}
		result.Message = fmt.Sprintf("Fetched and saved %d stock records", len(data))
	}

	result.Status = fly.StatusSuccess
	return result, nil
}

func (e *StockExecutor) saveToDatabase(ctx context.Context, collection string, data []bson.M) error {
	if e.db == nil {
		return fmt.Errorf("database not initialized")
	}

	coll, err := e.db.Collection(collection)
	if err != nil {
		return err
	}

	// 批量写入，使用 upsert
	var models []mongo.WriteModel
	for _, doc := range data {
		filter := bson.M{"code": doc["symbol"]}
		update := bson.M{"$set": doc}
		model := mongo.NewUpdateOneModel().SetFilter(filter).SetUpdate(update).SetUpsert(true)
		models = append(models, model)
	}

	opts := options.BulkWrite().SetOrdered(false)
	_, err = coll.BulkWrite(ctx, models, opts)
	return err
}
