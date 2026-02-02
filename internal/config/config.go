// Package config provides the application configuration loading functionality.
package config

// config package
import (
	"fmt"

	"fly-go/internal/database"

	"github.com/spf13/viper"
)

type Config struct {
	Server   ServerConfig
	Database database.Config
}

type ServerConfig struct {
	Port string
	Mode string
}

func LoadConfig() (*Config, error) {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("./config")

	viper.SetDefault("server.port", "8080")
	viper.SetDefault("server.mode", "debug")
	viper.SetDefault("database.host", "10.26.32.82")
	viper.SetDefault("database.port", "27017")
	viper.SetDefault("database.username", "root")
	viper.SetDefault("database.password", "123456")
	viper.SetDefault("database.database", "flygo")

	if err := viper.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			fmt.Println("Config file not found, using defaults")
		} else {
			return nil, fmt.Errorf("error reading config file: %w", err)
		}
	}

	viper.AutomaticEnv()

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		return nil, fmt.Errorf("error unmarshaling config: %w", err)
	}

	return &config, nil
}
