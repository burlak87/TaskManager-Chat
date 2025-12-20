package config

import (
	"sync"

	"github.com/your-team/taskmanager-chat/backend/pkg/logging"

	"github.com/ilyakaznacheev/cleanenv"
	"github.com/gin-gonic/gin"
    "github.com/spf13/viper"
)

// config/config.go
// package config

// import (
//     "github.com/gin-gonic/gin"
//     "github.com/spf13/viper"
// )

// type Config struct {
//     Port     string
//     GinMode  string
//     Database DatabaseConfig
// }

// func Load() *Config {
//     // Загрузка конфигурации
//     return &Config{
//         Port:    viper.GetString("PORT"),
//         GinMode: viper.GetString("GIN_MODE"),
//     }
// }

// func SetupRouter(config *Config) *gin.Engine {
//     if config.GinMode == "release" {
//         gin.SetMode(gin.ReleaseMode)
//     }
    
//     r := gin.New()
    
//     // Глобальные middleware
//     r.Use(gin.Logger())
//     r.Use(gin.Recovery())
    
//     return r
// }

type Config struct {
	Env           string `yml:"env" env-default:"development"`
	StorageConfig
}

type StorageConfig struct {
	Host     string `yaml:"host" env:"DB_HOST" env-default:"db"`
	Port     string `yaml:"port" env:"DB_PORT" env-default:"5432"`
	Database string `yaml:"database" env:"DB_NAME" env-default:"postgres"`
	Username string `yaml:"username" env:"DB_USER" env-default:"postgres"`
	Password string `yaml:"password" env:"DB_PASSWORD" env-default:"postgres"`
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
		
		logger.Infof("Database config: %s:%s", instance.Host, instance.Port)
	})
	return instance
}