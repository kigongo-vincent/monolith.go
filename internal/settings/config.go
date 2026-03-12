package settings

import (
	"os"

	"github.com/kigongo-vincent/monolith.go.git/internal/envloader"
)

// Config holds all app configuration. Secrets from .env.
type Config struct {
	Server  ServerConfig
	DB      DBConfig
	Storage StorageConfig
	Payment PaymentConfig
	Cache   CacheConfig
	Auth    AuthConfig
}

// ServerConfig for HTTP server.
type ServerConfig struct {
	Port string
}

// DBConfig for ORM (default SQLite).
type DBConfig struct {
	Driver string
	DSN    string
}

// StorageConfig: provider local|s3|cloudinary.
type StorageConfig struct {
	Provider  string
	LocalPath string
	S3        S3Config
	Cloud     CloudinaryConfig
}

type S3Config struct {
	Bucket string
	Region string
	Key    string
	Secret string
}

type CloudinaryConfig struct {
	CloudName string
	APIKey    string
	APISecret string
}

// PaymentConfig: provider stripe|flutterwave|pesapal.
type PaymentConfig struct {
	Provider string
}

// CacheConfig: backend redis|memory.
type CacheConfig struct {
	Backend string
	Redis   RedisConfig
}

type RedisConfig struct {
	URL string
}

// AuthConfig for Google SSO and JWT.
type AuthConfig struct {
	Google GoogleAuthConfig
	JWT    JWTConfig
}

type GoogleAuthConfig struct {
	ClientID     string
	ClientSecret string
}

type JWTConfig struct {
	Secret string
}

var global *Config

// Load reads .env from path then populates Config. Call from main. If path is empty, skips file.
func Load(envPath string) (*Config, error) {
	if envPath != "" {
		if err := envloader.Load(envPath); err != nil {
			return nil, err
		}
	}
	c := &Config{
		Server:  ServerConfig{Port: getEnv("PORT", "8080")},
		DB:      DBConfig{Driver: getEnv("DB_DRIVER", "sqlite"), DSN: getEnv("DB_DSN", "file:monolith.db")},
		Storage: StorageConfig{Provider: getEnv("STORAGE_PROVIDER", "local"), LocalPath: getEnv("STORAGE_LOCAL_PATH", "storage")},
		Payment: PaymentConfig{Provider: getEnv("PAYMENT_PROVIDER", "stripe")},
		Cache:   CacheConfig{Backend: getEnv("CACHE_BACKEND", "memory"), Redis: RedisConfig{URL: getEnv("REDIS_URL", "")}},
		Auth: AuthConfig{
			Google: GoogleAuthConfig{ClientID: getEnv("GOOGLE_CLIENT_ID", ""), ClientSecret: getEnv("GOOGLE_CLIENT_SECRET", "")},
			JWT:    JWTConfig{Secret: getEnv("JWT_SECRET", "change-me")},
		},
	}
	global = c
	return c, nil
}

// Get returns the global config (nil until Load is called).
func Get() *Config { return global }

func getEnv(key, def string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return def
}
