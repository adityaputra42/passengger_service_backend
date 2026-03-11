package config

import (
	"log"
	"strings"
	"time"

	"github.com/spf13/viper"
)

type Config struct {
	Database DatabaseConfig
	Server   ServerConfig
	JWT      JWTConfig
	CORS     CORSConfig
	SMTP     SMTPConfig
	System   SystemConfig
	Supabase SupabaseConfig
}

type SupabaseConfig struct {
	Url string
	Key string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

type ServerConfig struct {
	Port string
	Host string
	Env  string
}

type JWTConfig struct {
	Secret             string
	RefreshSecret      string
	AccessTokenExpiry  time.Duration
	RefreshTokenExpiry time.Duration
}

type CORSConfig struct {
	AllowedOrigins []string
}

type SMTPConfig struct {
	Host     string
	Port     string
	Username string
	Password string
}

type SystemConfig struct {
	DefaultAdminEmail    string
	DefaultAdminPassword string
}

func Load() *Config {
	viper.SetConfigFile(".env")
	viper.AutomaticEnv()

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("Warning: Could not load .env file: %v", err)
	}

	accessTokenExpiry, _ := time.ParseDuration(viper.GetString("JWT_ACCESS_TOKEN_EXPIRY"))
	refreshTokenExpiry, _ := time.ParseDuration(viper.GetString("JWT_REFRESH_TOKEN_EXPIRY"))

	allowedOrigins := strings.Split(viper.GetString("CORS_ALLOWED_ORIGINS"), ",")
	for i := range allowedOrigins {
		allowedOrigins[i] = strings.TrimSpace(allowedOrigins[i])
	}

	return &Config{
		Database: DatabaseConfig{
			Host:     viper.GetString("DB_HOST"),
			Port:     viper.GetString("DB_PORT"),
			User:     viper.GetString("DB_USER"),
			Password: viper.GetString("DB_PASSWORD"),
			Name:     viper.GetString("DB_NAME"),
			SSLMode:  viper.GetString("DB_SSL_MODE"),
		},
		Server: ServerConfig{
			Port: viper.GetString("PORT"),
			Host: viper.GetString("HOST"),
			Env:  viper.GetString("ENV"),
		},
		JWT: JWTConfig{
			Secret:             viper.GetString("JWT_SECRET"),
			RefreshSecret:      viper.GetString("JWT_REFRESH_SECRET"),
			AccessTokenExpiry:  accessTokenExpiry,
			RefreshTokenExpiry: refreshTokenExpiry,
		},
		CORS: CORSConfig{
			AllowedOrigins: allowedOrigins,
		},
		SMTP: SMTPConfig{
			Host:     viper.GetString("SMTP_HOST"),
			Port:     viper.GetString("SMTP_PORT"),
			Username: viper.GetString("SMTP_USERNAME"),
			Password: viper.GetString("SMTP_PASSWORD"),
		},
		System: SystemConfig{
			DefaultAdminEmail:    viper.GetString("DEFAULT_ADMIN_EMAIL"),
			DefaultAdminPassword: viper.GetString("DEFAULT_ADMIN_PASSWORD"),
		},
		Supabase: SupabaseConfig{
			Url: viper.GetString("SUPABASE_URL"),
			Key: viper.GetString("SUPABASE_KEY"),
		},
	}
}
