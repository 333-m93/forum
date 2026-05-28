package config

import "os"

type Config struct {
	Port   string
	DBUser string
	DBPass string
	DBHost string
	DBName string
}

func Load() Config {
	return Config{
		Port:   getEnv("APP_PORT", ":8080"),
		DBUser: getEnv("DB_USER", "root"),
		DBPass: getEnv("DB_PASS", ""),
		DBHost: getEnv("DB_HOST", "127.0.0.1:3306"),
		DBName: getEnv("DB_NAME", "forumdb"),
	}
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
