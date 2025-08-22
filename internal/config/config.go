package config

import (
	"github.com/spf13/viper"
	"log"
	"time"
)

// Config holds the application's configuration
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	JWT      JWTConfig
	Redis    RedisConfig
	OSS      OSSConfig
}

// ServerConfig holds server settings
type ServerConfig struct {
	Port string
	Mode string
}

// DatabaseConfig holds database settings
type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	DBName   string
	SSLMode  string
}

// JWTConfig holds JWT settings
type JWTConfig struct {
	Secret      string
	Issuer      string
	ExpireHours time.Duration
}

// RedisConfig holds Redis settings
type RedisConfig struct {
	Addr     string
	Password string
	DB       int
}

// OSSConfig holds Alibaba Cloud OSS settings
type OSSConfig struct {
	Endpoint        string
	AccessKeyID     string
	AccessKeySecret string
	BucketName      string
}

var AppConfig *Config

// LoadConfig loads configuration from file and environment variables
func LoadConfig(path string) {
	viper.AddConfigPath(path)
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")

	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Failed to read config file: %v", err)
	}

	var config Config
	if err := viper.Unmarshal(&config); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	// Post-process JWT expiration
	config.JWT.ExpireHours = time.Duration(viper.GetInt("jwt.expire_hours")) * time.Hour

	AppConfig = &config
	log.Println("Configuration loaded successfully.")
}
