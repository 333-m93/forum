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
	cfg := Config{
		Port:   getEnv("APP_PORT", ":8080"),
		DBUser: getEnv("DB_USER", ""),
		DBPass: getEnv("DB_PASS", ""),
		DBHost: getEnv("DB_HOST", ""),
		DBName: getEnv("DB_NAME", ""),
	}

	if cfg.DBHost == "" || cfg.DBUser == "" || cfg.DBName == "" {
		panic("DB configuration missing (check environment variables)")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value == "" {
		return fallback
	}
	return value
}
