package config

import (
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	Name     string
	User     string
	Password string
	SSLMode  string
}

// RedisConfig holds configuration for connecting to a Redis server.
// Redis is an in-memory NoSQL data structure store, used as a database, cache, and message broker.
type RedisConfig struct {
	Host     string
	Port     string
	Password string
	DB       int
}

type JWTConfig struct {
	Secret            string
	Expiration        time.Duration
	RefreshExpiration time.Duration
}

var AppConfig *Config

func LoadConfig() (*Config, error) {
	// godotenv is a library for loading environment variables
	_ = godotenv.Load()

	AppConfig = &Config{
		Server: ServerConfig{
			Port: getEnv("APP_PORT"),
			Host: getEnv("APP_HOST"),
			Env:  getEnv("APP_ENV"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST"),
			Port:     getEnv("DB_PORT"),
			Name:     getEnv("DB_NAME"),
			User:     getEnv("DB_USER"),
			Password: getEnv("DB_PASSWORD"),
			SSLMode:  getEnv("DB_SSLMODE"),
		},
		Redis: RedisConfig{
			Host:     getEnv("REDIS_HOST"),
			Port:     getEnv("REDIS_PORT"),
			Password: getEnv("REDIS_PASSWORD"),
			DB:       getEnvInt("REDIS_DB"),
		},
		JWT: JWTConfig{
			Secret:            getEnv("JWT_SECRET"),
			Expiration:        getEnvDuration("JWT_EXPIRATION"),
			RefreshExpiration: getEnvDuration("JWT_REFRESH_EXPIRATION"),
		},
	}

	if err := validateConfig(AppConfig); err != nil {
		return nil, err
	}
	return AppConfig, nil
}

func (c *Config) GetDatabaseURL() string {
    return fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=%s",
        c.Database.Host,
        c.Database.Port,
        c.Database.User,
        c.Database.Password,
        c.Database.Name,
        c.Database.SSLMode,
    )
}

func (c *Config) GetRedisURL() string {
    if c.Redis.Password != "" {
        return fmt.Sprintf("redis://:%s@%s:%s/%d",
            c.Redis.Password,
            c.Redis.Host,
            c.Redis.Port,
            c.Redis.DB,
        )
    }
    return fmt.Sprintf("redis://%s:%s/%d",
        c.Redis.Host,
        c.Redis.Port,
        c.Redis.DB,
    )
}

func getEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		panic(fmt.Sprintf("Required environment variable not set: %s", key))
	}
	return val
}

func getEnvInt(key string) int {
	val := getEnv(key)
	i, err := strconv.Atoi(val)
	if err != nil {
		panic(fmt.Sprintf("Variable %s must be a valid integer", key))
	}
	return i
}

func getEnvDuration(key string) time.Duration {
	val := getEnv(key)
	d, err := time.ParseDuration(val)
	if err != nil {
		panic(fmt.Sprintf("Variable %s must be a valid duration", key))
	}
	return d
}

func validateConfig(config *Config) error {
	if config.JWT.Secret == "" {
		return fmt.Errorf("JWT_SECRET is required")
	}
	if config.Database.Password == "" {
		return fmt.Errorf("DB_PASSWORD is required")
	}
	return nil
}
