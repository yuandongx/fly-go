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
	config Config
}

type Row = map[string]interface{}
type Rows = []Row

func NewMongoDB(config Config) (*MongoDB, error) {
	m := &MongoDB{}
	m.config = config
	err := m.Connect()
	return m, err
}

func (m *MongoDB) Connect() error {
	clientOptions := options.Client().ApplyURI("mongodb://" + m.config.Host + ":" + m.config.Port)
	clientOptions.SetAuth(options.Credential{
		Username: m.config.Username,
		Password: m.config.Password,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return err
	}
	m.Client = client
	m.DB = client.Database(m.config.Database)
	return nil
}

// Close disconnects the MongoDB client.
func (m *MongoDB) Close() error {
	return m.Client.Disconnect(context.Background())
}

func (m *MongoDB) Collection(name string) *mongo.Collection {
	ctxt := context.Background()
	if err := m.Client.Ping(ctxt, nil); err != nil {
		if err := m.Connect(); err != nil {
			return nil
		}
	}
	return m.DB.Collection(name)
}
