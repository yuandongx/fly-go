// Package database provides the MongoDB connection and interaction functionality.
package database

// database package
import (
	"context"
	"time"

	"fly-go/logger"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

var Client *mongo.Client
var DB *mongo.Database

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

func Connect(config Config) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	uri := "mongodb://" + config.Host + ":" + config.Port
	clientOptions := options.Client().ApplyURI(uri)
	clientOptions.SetAuth(options.Credential{
		Username: config.Username,
		Password: config.Password,
	})

	var err error
	Client, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		logger.Error("Failed to connect to MongoDB", zap.Error(err))
		return err
	}

	if err = Client.Ping(ctx, nil); err != nil {
		logger.Error("Failed to ping MongoDB", zap.Error(err))
		return err
	}

	DB = Client.Database(config.Database)

	logger.Info("Successfully connected to MongoDB",
		zap.String("database", config.Database))

	return nil
}

func Disconnect() error {
	if Client == nil {
		return nil
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := Client.Disconnect(ctx); err != nil {
		logger.Error("Failed to disconnect from MongoDB", zap.Error(err))
		return err
	}

	logger.Info("Successfully disconnected from MongoDB")
	return nil
}

func GetCollection(name string) *mongo.Collection {
	return DB.Collection(name)
}
