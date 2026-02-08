// Package config load config from yaml
package config

import (
	"fmt"
	"io"
	"os"
	"reflect"
	"strconv"

	"gopkg.in/yaml.v3"
)

type DatabaseConfig struct {
	MongoHost     string `yaml:"mongo_host" env:"MONGO_HOST" default:"localhost"`
	MongoPort     int    `yaml:"mongo_port" env:"MONGO_PORT" default:"27017"`
	MongoPassword string `yaml:"mongo_password" env:"MONGO_PASSWORD" default:"mypassword"`
	MongoDatabase string `yaml:"mongo_database" env:"MONGO_DATABASE" default:"mydb"`
	MongoUsername string `yaml:"mongo_username" env:"MONGO_USERNAME" default:"myuser"`
}

type ServerConfig struct {
	Port int    `yaml:"port" env:"PORT" default:"8080"`
	Mode string `yaml:"mode" env:"MODE" default:"release"`
}

type Config struct {
	Name     string         `yaml:"name"`
	Database DatabaseConfig `yaml:"database"`
	Server   ServerConfig   `yaml:"server"`
}

// LoadConfig 加载并返回应用程序配置，包含默认值和从YAML文件加载的配置。
// 返回的配置结构包含应用名称、数据库配置和服务器配置(默认端口8080，模式为release)。
// 如果加载YAML文件失败，返回的错误将包含具体原因。
func LoadConfig() (*Config, error) {
	cfg := &Config{
		Name:     "fly-go",
		Database: DatabaseConfig{},
		Server:   ServerConfig{Port: 8000, Mode: "release"},
	}
	file, ok := os.LookupEnv("FLY_CONFIG")
	if !ok {
		file = "config.yaml"
	}
	err := LoadYAML(file, cfg)
	return cfg, err
}

// LoadYAML 从YAML配置文件加载配置到给定的Config结构体
// 默认读取config.yaml文件，可通过FLY_CONFIG环境变量指定自定义配置文件路径
// 返回错误如果文件不存在、读取失败或解析失败
func LoadYAML(file string, cfg *Config) error {
	if _, err := os.Stat(file); os.IsNotExist(err) {
		return err
	}
	r, err := os.Open(file)
	if err != nil {
		return err
	}
	defer r.Close()
	if data, err := io.ReadAll(r); err != nil {
		return err
	} else {
		return yaml.Unmarshal(data, cfg)
	}

}

// LoadEnv 从环境变量加载配置并覆盖对应的配置项
// 环境变量优先级高于配置文件
func LoadEnv(cfg *Config, force bool) error {
	// 从环境变量加载配置并覆盖对应配置项，支持嵌套结构体
	v := reflect.ValueOf(cfg)
	if v.Kind() != reflect.Ptr || v.IsNil() {
		return fmt.Errorf("cfg must be a non-nil pointer")
	}

	var load func(reflect.Value) error
	load = func(rv reflect.Value) error {
		if rv.Kind() == reflect.Ptr {
			rv = rv.Elem()
		}
		if rv.Kind() != reflect.Struct {
			return load(rv)
		}
		for i := 0; i < rv.NumField(); i++ {
			field := rv.Field(i)
			ftype := rv.Type().Field(i)
			fmt.Printf("Process field: %s\n", ftype.Name)
			// If the field is a nested struct, recurse
			if field.Kind() == reflect.Struct {
				if err := load(field); err != nil {
					return err
				}
			}

			tag := ftype.Tag.Get("env")
			defaultValue := ftype.Tag.Get("default")

			// 如果环境变量不存在且字段值为零值，则使用默认值
			env, exists := os.LookupEnv(tag)
			if !exists {
				if field.IsNil() {
					// 如果字段值为零值且存在默认值，则使用默认值
					if defaultValue == "" {
						continue
					}
					env = defaultValue
				} else {
					continue
				}
			}

			// 根据字段类型设置值，支持字符串、整数和布尔类型
			switch field.Kind() {
			case reflect.String:
				field.SetString(env)
			case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
				if env != "" {
					if iv, err := strconv.ParseInt(env, 10, 64); err == nil {
						field.SetInt(iv)
						continue
					}
				}
				if defaultValue != "" {
					if dv, err := strconv.ParseInt(defaultValue, 10, 64); err == nil {
						field.SetInt(dv)
					}
				}
			case reflect.Bool:
				if env != "" {
					if bv, err := strconv.ParseBool(env); err == nil {
						field.SetBool(bv)
						continue
					}
				}
				if defaultValue != "" {
					if db, err := strconv.ParseBool(defaultValue); err == nil {
						field.SetBool(db)
					}
				}
			default:
				return fmt.Errorf("unsupported field type: %s", field.Kind())
			}
		}
		return nil
	}

	return load(v)
}
