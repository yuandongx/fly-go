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
	mg := &MongoDB{}
	mg.config = config
	error := mg.Connect()
	return mg, error
}
func (mg *MongoDB) Connect() error {
	clientOptions := options.Client().ApplyURI("mongodb://" + mg.config.Host + ":" + mg.config.Port)
	clientOptions.SetAuth(options.Credential{
		Username: mg.config.Username,
		Password: mg.config.Password,
	})

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return err
	}

	if err = client.Ping(context.Background(), nil); err != nil {
		return err
	}
	mg.Client = client
	mg.DB = client.Database(mg.config.Database)
	return nil
}

func (mg *MongoDB) Close() error {
	return mg.Client.Disconnect(context.Background())
}

func (mg *MongoDB) Collection(name string) *mongo.Collection {
	ctxt := context.Background()
	if err := mg.Client.Ping(ctxt, nil); err != nil {
		if error := mg.Connect(); error != nil {
			return nil
		}
	}
	return mg.DB.Collection(name)
}
