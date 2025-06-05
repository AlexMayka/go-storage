package config

import "os"

type Minio struct {
	User     string
	Password string
	Port     string
	Host     string
}

type Db struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
}

type App struct {
	Host      string
	Port      string
	JwtSecret string
	LogToFile string
	LogLevel  string
}

type Config struct {
	Minio Minio
	Db    Db
	App   App
}

func NewConfig() *Config {
	return &Config{
		Minio: Minio{
			User:     GetEnv("MINIO_ROOT_USER", "admin"),
			Password: GetEnv("MINIO_ROOT_PASSWORD", "secret123"),
			Port:     GetEnv("MINIO_API_PORT", "9000"),
			Host:     GetEnv("MINIO_ROOT_HOST", "localhost"),
		},
		Db: Db{
			Host:     GetEnv("POSTGRES_HOST", "localhost"),
			Port:     GetEnv("POSTGRES_PORT", "5432"),
			User:     GetEnv("POSTGRES_USER", "admin"),
			Password: GetEnv("POSTGRES_PASSWORD", "admin"),
			Name:     GetEnv("POSTGRES_DB", "storage"),
		},
		App: App{
			Host:      GetEnv("APP_HOST", "localhost"),
			Port:      GetEnv("APP_PORT", "8080"),
			JwtSecret: GetEnv("APP_JWT_SECRET", "secret"),
			LogToFile: GetEnv("APP_LOG_FILE", "true"),
			LogLevel:  GetEnv("APP_LOG_LEVEL", "info"),
		},
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
