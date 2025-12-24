package config

import (
	"sync"

	"github.com/your-team/taskmanager-chat/backend/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env string `yml:"env" env-default:"development"`
	StorageConfig
	MongoConfig
}

type StorageConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"db"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Database string `yaml:"database" env:"DB_NAME" env-default:"postgres"`
	Username string `yaml:"username" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
}

type MongoConfig struct {
	Host     string `yaml:"host" env:"MONGO_HOST" env-default:"mongodb"`
	Port     string `yaml:"port" env:"MONGO_PORT" env-default:"27017"`
	Database string `yaml:"database" env:"MONGO_DB" env-default:"taskmanager"`
	Username string `yaml:"username" env:"MONGO_USER" env-default:"admin"`
	Password string `yaml:"password" env:"MONGO_PASSWORD" env-default:"password"`
}

var instance *Config
var once sync.Once

func GetConfig() *Config {
	once.Do(func() {
		logger := logging.GetLogger()
		logger.Info("read application configuration")
		instance = &Config{}

		if err := cleanenv.ReadEnv(instance); err != nil {
			logger.Errorf("Error reading env vars: %v", err)
		}

		if err := cleanenv.ReadConfig("/config.yml", instance); err != nil {
			logger.Warnf("Config file not found, using env vars: %v", err)
		}

		logger.Infof("Database config: %s:%s", instance.StorageConfig.Host, instance.StorageConfig.Port)
	})
	return instance
}
