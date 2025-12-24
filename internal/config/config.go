package config

import (
	"fmt"
	"os"

	"subscribe_project/pkg/logger"

	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"
)

type Config struct {
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	ServerPort string
}

func LoadConfig() (*Config, error) {
	logger.Log.Info("Loading configuration...")

	if err := godotenv.Load(".env"); err != nil {
		logger.Log.WithField("error", err).Warn(".env file not found, using environment variables")
	}

	config := &Config{
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "postgres"),
		DBName:     getEnv("DB_NAME", "subscriptions_db"),
		ServerPort: getEnv("SERVER_PORT", "8080"),
	}

	logger.Log.WithFields(logrus.Fields{
		"db_host":     config.DBHost,
		"db_port":     config.DBPort,
		"db_name":     config.DBName,
		"server_port": config.ServerPort,
	}).Info("Configuration loaded successfully")

	return config, nil
}

func (c *Config) GetDBConnectionString() string {
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		c.DBHost, c.DBPort, c.DBUser, c.DBPassword, c.DBName)

	logger.Log.WithFields(logrus.Fields{
		"host":   c.DBHost,
		"port":   c.DBPort,
		"user":   c.DBUser,
		"dbname": c.DBName,
	}).Debug("Generated database connection string")

	return connStr
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {

		if key == "DB_PASSWORD" {
			logger.Log.WithField(key, "***").Debug("Loaded environment variable")
		} else {
			logger.Log.WithField(key, value).Debug("Loaded environment variable")
		}
		return value
	}

	logger.Log.WithFields(logrus.Fields{
		"key":           key,
		"default_value": defaultValue,
	}).Debug("Using default value for environment variable")

	return defaultValue
}
