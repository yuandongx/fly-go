// Package config load config from yaml
package config

type Config struct {
	Name          string `yaml:"name"`
	MongoHost     string `yaml:"mongo_host"`
	MongoPort     string `yaml:"mongo_port"`
	MongoPassword string `yaml:"mongo_password"`
	MongoDatabase string `yaml:"mongo_database"`
}
