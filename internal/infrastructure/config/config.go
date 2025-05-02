package config

import (
	"fmt"
	"os"
	"strconv"
	"strings"
)

type Config struct {
	Port          string
	MinioEndPoint string
	BucketName    string
	MinioUser     string
	MinioPassword string
	MinioSSL      bool
}

func LoadConfig() (*Config, error) {
	appConfig := &Config{
		Port: getEnv("PORT", "8080"),

		MinioEndPoint: getEnv("MINIO_ENDPOINT", "localhost:9000"),
		BucketName:    getEnv("MINIO_BUCKET_NAME", "defaultBucket"),
		MinioUser:     getEnv("MINIO_ROOT_USER", "root"),
		MinioPassword: getEnv("MINIO_ROOT_PASSWORD", "minio_password"),
		MinioSSL:      getEnvAsBool("MINIO_USE_SSL", false),
	}
	if err := validateConfig(appConfig); err != nil {
		return nil, err
	}
	return appConfig, nil
}

func validateConfig(cfg *Config) error {
	if cfg.Port == "" {
		return fmt.Errorf("port is required")
	}
	if !strings.HasPrefix(cfg.Port, ":") {
		cfg.Port = ":" + cfg.Port
	}

	if cfg.MinioEndPoint == "" {
		return fmt.Errorf("MinIO endpoint is required")
	}
	if cfg.BucketName == "" {
		return fmt.Errorf("bucket name is required")
	}
	if cfg.MinioUser == "" {
		return fmt.Errorf("MinIO user is required")
	}
	if cfg.MinioPassword == "" {
		return fmt.Errorf("MinIO password is required")
	}

	return nil
}

func getEnv(key string, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if valueStr := getEnv(key, ""); valueStr != "" {
		if value, err := strconv.ParseBool(valueStr); err == nil {
			return value
		}
	}

	return defaultValue
}
