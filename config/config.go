package config

import (
	"os"
)

func getEnv(key string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return ""
}

type PropertyConfig struct {
	Database struct {
		Host     string
		Port     string
		DBName   string
		UserName string
		Password string
		SSLMode  string
	}
	App struct {
		GiteaBaseURL  string
		GiteaToken    string
		GiteaUsername string
	}
}

var PropConfig *PropertyConfig = LoadConfigFromEnv()

func LoadConfigFromEnv() *PropertyConfig {
	cfg := &PropertyConfig{}

	cfg.Database.Host = getEnv("DB_HOST")
	cfg.Database.Port = getEnv("DB_PORT")
	cfg.Database.DBName = getEnv("DB_NAME")
	cfg.Database.UserName = getEnv("DB_USER")
	cfg.Database.Password = getEnv("DB_PASSWORD")
	cfg.Database.SSLMode = getEnv("DB_SSLMODE")

	cfg.App.GiteaBaseURL = getEnv("GITEA_BASE_URL")
	cfg.App.GiteaToken = getEnv("GITEA_TOKEN")
	cfg.App.GiteaUsername = getEnv("GITEA_USERNAME")
	return cfg
}
