package config

import (
	"os"
	"strconv"
	"time"
)

type Minio struct {
	User       string
	Password   string
	Port       string
	Host       string
	BucketName string
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

type FileServer struct {
	SmallFileThreshold  int64
	MediumFileThreshold int64
	LargeFileThreshold  int64
	MaxFileSize         int64

	MaxConcurrentUploads   int
	MaxMemoryPerRequest    int64
	MaxTotalMemoryForFiles int64

	ChunkSize          int64
	ChunkUploadTimeout time.Duration
	ChunkedSessionTTL  time.Duration

	BufferSize      int
	DownloadTimeout time.Duration

	MemoryPressureThreshold float64
	CPUPressureThreshold    float64

	MaxFailuresBeforeOpen int
	CircuitBreakerTimeout time.Duration
}

type Config struct {
	Minio      Minio
	Db         Db
	App        App
	FileServer FileServer
}

func NewConfig() *Config {
	return &Config{
		Minio: Minio{
			User:       GetEnv("MINIO_ROOT_USER", "admin"),
			Password:   GetEnv("MINIO_ROOT_PASSWORD", "secret123"),
			Port:       GetEnv("MINIO_API_PORT", "9000"),
			Host:       GetEnv("MINIO_ROOT_HOST", "localhost"),
			BucketName: GetEnv("MINIO_BUCKET_NAME", "go-storage"),
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
		FileServer: FileServer{
			SmallFileThreshold:  GetEnvInt64("FILE_SMALL_THRESHOLD", 10*1024*1024),
			MediumFileThreshold: GetEnvInt64("FILE_MEDIUM_THRESHOLD", 100*1024*1024),
			LargeFileThreshold:  GetEnvInt64("FILE_LARGE_THRESHOLD", 1024*1024*1024),
			MaxFileSize:         GetEnvInt64("FILE_MAX_SIZE", 5*1024*1024*1024),

			MaxConcurrentUploads:   GetEnvInt("FILE_MAX_CONCURRENT_UPLOADS", 10),
			MaxMemoryPerRequest:    GetEnvInt64("FILE_MAX_MEMORY_PER_REQUEST", 100*1024*1024),
			MaxTotalMemoryForFiles: GetEnvInt64("FILE_MAX_TOTAL_MEMORY", 500*1024*1024),

			ChunkSize:          GetEnvInt64("FILE_CHUNK_SIZE", 5*1024*1024),
			ChunkUploadTimeout: GetEnvDuration("FILE_CHUNK_TIMEOUT", 30*time.Second),
			ChunkedSessionTTL:  GetEnvDuration("FILE_CHUNK_SESSION_TTL", 24*time.Hour),

			BufferSize:      GetEnvInt("FILE_BUFFER_SIZE", 64*1024),
			DownloadTimeout: GetEnvDuration("FILE_DOWNLOAD_TIMEOUT", 10*time.Minute),

			MemoryPressureThreshold: GetEnvFloat64("FILE_MEMORY_PRESSURE_THRESHOLD", 0.8),
			CPUPressureThreshold:    GetEnvFloat64("FILE_CPU_PRESSURE_THRESHOLD", 0.7),

			MaxFailuresBeforeOpen: GetEnvInt("FILE_CIRCUIT_MAX_FAILURES", 5),
			CircuitBreakerTimeout: GetEnvDuration("FILE_CIRCUIT_TIMEOUT", 1*time.Minute),
		},
	}
}

func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

func GetEnvInt64(key string, fallback int64) int64 {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseInt(value, 10, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func GetEnvInt(key string, fallback int) int {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.Atoi(value); err == nil {
			return parsed
		}
	}
	return fallback
}

func GetEnvFloat64(key string, fallback float64) float64 {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := strconv.ParseFloat(value, 64); err == nil {
			return parsed
		}
	}
	return fallback
}

func GetEnvDuration(key string, fallback time.Duration) time.Duration {
	if value, ok := os.LookupEnv(key); ok {
		if parsed, err := time.ParseDuration(value); err == nil {
			return parsed
		}
	}
	return fallback
}
