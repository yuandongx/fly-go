// Package database provides the MongoDB connection and interaction functionality.
package database

// database package
import (
	"context"
	"errors"
	"fly-go/config"
	"fmt"
	"strconv"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type MongoDB struct {
	Client *mongo.Client
	DB     *mongo.Database
	config config.DatabaseConfig
}

type Row = map[string]interface{}
type Rows = []Row

func NewMongoDB(cfg config.DatabaseConfig) (*MongoDB, error) {
	m := &MongoDB{}
	m.config = cfg
	err := m.Connect()
	return m, err
}

func (m *MongoDB) Connect() error {
	// Configure connection pool with reasonable defaults for production
	fmt.Printf("Connecting to MongoDB at %s:%d...\n", m.config.MongoHost, m.config.MongoPort)
	clientOptions := options.Client().
		ApplyURI("mongodb://" + m.config.MongoHost + ":" + strconv.Itoa(m.config.MongoPort)).
		SetMaxPoolSize(100).                        // Maximum number of connections
		SetMinPoolSize(10).                         // Minimum number of connections
		SetMaxConnIdleTime(30 * time.Second).       // Close idle connections after 30s
		SetConnectTimeout(10 * time.Second).        // Timeout for initial connection
		SetServerSelectionTimeout(10 * time.Second) // Timeout for server selection

	clientOptions.SetAuth(options.Credential{
		Username: m.config.MongoUsername,
		Password: m.config.MongoPassword,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		_ = client.Disconnect(context.Background())
		return err
	}
	m.Client = client
	m.DB = client.Database(m.config.MongoDatabase)
	return nil
}

// Close disconnects the MongoDB client.
func (m *MongoDB) Close() error {
	return m.Client.Disconnect(context.Background())
}

func (m *MongoDB) Collection(name string) (*mongo.Collection, error) {
	if m.Client == nil {
		if err := m.Connect(); err != nil {
			return nil, err
		}
	}

	ctxt, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := m.Client.Ping(ctxt, nil); err != nil {
		_ = m.Close()
		if err := m.Connect(); err != nil {
			return nil, errors.New("failed to reconnect to MongoDB: " + err.Error())
		}
	}
	return m.DB.Collection(name), nil
}
