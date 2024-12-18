package config

import (
	"os"
)

type Configuration struct {
	DatabaseName        string
	DatabaseHost        string
	DatabasePort        string
	DatabaseUser        string
	DatabasePassword    string
	MigrateToVersion    string
	MigrationLocation   string
	CloudinaryNameKey   string
	CloudinaryApiKey    string
	CloudinarySecretKey string
}

func GetConfiguration() Configuration {
	return Configuration{
		DatabaseName:        getOrDefault("DB_NAME", "restapi_dev"),
		DatabaseHost:        getOrDefault("DB_HOST", "127.0.0.1"),
		DatabasePort:        getOrDefault("DB_PORT", "5432"),
		DatabaseUser:        getOrDefault("DB_USER", "postgres"),
		DatabasePassword:    getOrDefault("DB_PASSWORD", "postgres"),
		MigrateToVersion:    getOrDefault("MIGRATE", "latest"),
		MigrationLocation:   getOrDefault("MIGRATION_LOCATION", "migrations"),
		CloudinaryNameKey:   os.Getenv("CLOUDINARY_NAME_KEY"),
		CloudinaryApiKey:    os.Getenv("CLOUDINARY_API_KEY"),
		CloudinarySecretKey: os.Getenv("CLOUDINARY_SECRET_KEY"),
	}
}

func getOrDefault(key, defaultVal string) string {
	env, set := os.LookupEnv(key)
	if !set {
		return defaultVal
	}
	return env
}
