package config

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	Env  string
	Port string
}

type DatabaseConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	TimeZone string
}

type JWTConfig struct {
	Secret          string
	ExpirationHours int
}

type WireGuardConfig struct {
	Subnet     string
	DNS        string
	AllowedIPs string
	Keepalive  int
}

type Config struct {
	App       AppConfig
	Database  DatabaseConfig
	JWT       JWTConfig
	WireGuard WireGuardConfig
}

func Load() *Config {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env file found, falling back to system environment")
	}

	return &Config{
		App: AppConfig{
			Env:  getEnv("APP_ENV", "development"),
			Port: getEnv("APP_PORT", "8080"),
		},
		Database: DatabaseConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "postgres"),
			Password: getEnv("DB_PASSWORD", "postgres"),
			Name:     getEnv("DB_NAME", "wg-vpn-platform"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
			TimeZone: getEnv("DB_TIMEZONE", "UTC"),
		},
		JWT: JWTConfig{
			Secret:          getEnv("JWT_SECRET", "change_me"),
			ExpirationHours: getEnvAsInt("JWT_EXPIRATION_HOURS", 2160), // 90 days (~3 months)
		},
		WireGuard: WireGuardConfig{
			Subnet:     getEnv("WG_DEFAULT_SUBNET", "10.8.0.0/24"),
			DNS:        getEnv("WG_DNS", "1.1.1.1,8.8.8.8"),
			AllowedIPs: getEnv("WG_ALLOWED_IPS", "0.0.0.0/0,::/0"),
			Keepalive:  getEnvAsInt("WG_KEEPALIVE", 25),
		},
	}
}

func getEnv(key, fallback string) string {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		return v
	}
	return fallback
}

func getEnvAsInt(key string, fallback int) int {
	if v, ok := os.LookupEnv(key); ok && v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			return n
		}
	}
	return fallback
}

func (d DatabaseConfig) DSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		d.Host, d.Port, d.User, d.Password, d.Name, d.SSLMode, d.TimeZone,
	)
}
