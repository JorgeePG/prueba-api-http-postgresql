package config

import (
	"fmt"
	"os"
)

type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
}

type DatabaseConfig struct {
	Host     string
	Port     int
	User     string
	Password string
	DBName   string
}

type ServerConfig struct {
	Port    string
	SSLCert string
	SSLKey  string
	UseSSL  bool
}

func Load() *Config {
	return &Config{
		Server: ServerConfig{
			Port:    getEnv("SERVER_PORT", "8080"),
			SSLCert: getEnv("SSL_CERT_PATH", "certs/ssl/server.crt"),
			SSLKey:  getEnv("SSL_KEY_PATH", "certs/ssl/server.key"),
			UseSSL:  getEnv("USE_SSL", "false") == "true",
		},
		Database: DatabaseConfig{
			Host:     "localhost",
			Port:     5432,
			DBName:   "api_db",
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
		},
	}
}

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func (d *DatabaseConfig) ConnectionString() string {

	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		d.Host, d.Port, d.User, d.Password, d.DBName)
}
