// Package database provides the MongoDB connection and interaction functionality.
package database

// database package
import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Config struct {
	Host     string
	Port     string
	Username string
	Password string
	Database string
}

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
	Config Config
}

type Row = map[string]interface{}
type Rows = []Row

func NewMongoDB(config Config) (*MongoDB, error) {
	clientOptions := options.Client().ApplyURI("mongodb://" + config.Host + ":" + config.Port)
	clientOptions.SetAuth(options.Credential{
		Username: config.Username,
		Password: config.Password,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return nil, err
	}

	return &MongoDB{
		Client: client,
		DB:     client.Database(config.Database),
		Config: config,
	}, nil
}

func (m *MongoDB) Close() error {
	return m.Client.Disconnect(context.Background())
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	return m.DB.Collection(name)
}
