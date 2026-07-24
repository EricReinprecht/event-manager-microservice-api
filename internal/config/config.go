package config

import (
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret string

	TicketVerificationTTL time.Duration

	PayPalClientID     string
	PayPalClientSecret string
	PayPalBaseURL      string
}

func Load() *Config {

	godotenv.Load(".env")
	godotenv.Overload(".env.local")

	return &Config{

		Port: getEnv(
			"PORT",
			"8080",
		),

		Env: getEnv(
			"APP_ENV",
			"development",
		),

		DBHost: getEnv(
			"DB_HOST",
			"localhost",
		),

		DBPort: getEnv(
			"DB_PORT",
			"5432",
		),

		DBUser: getEnv(
			"DB_USER",
			"postgres",
		),

		DBPassword: getEnv(
			"DB_PASSWORD",
			"",
		),

		DBName: getEnv(
			"DB_NAME",
			"event_platform",
		),

		JWTSecret: getEnv(
			"JWT_SECRET",
			"development-secret-change-me",
		),

		TicketVerificationTTL: time.Duration(
			getEnvInt(
				"TICKET_VERIFICATION_TTL_MINUTES",
				15,
			),
		) * time.Minute,

		PayPalClientID: getEnv(
			"PAYPAL_CLIENT_ID",
			"",
		),

		PayPalClientSecret: getEnv(
			"PAYPAL_CLIENT_SECRET",
			"",
		),

		PayPalBaseURL: getEnv(
			"PAYPAL_BASE_URL",
			"https://api-m.sandbox.paypal.com",
		),
	}
}

func getEnv(key string, fallback string) string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	return value
}

func getEnvInt(
	key string,
	fallback int,
) int {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	result, err := strconv.Atoi(value)

	if err != nil {
		return fallback
	}

	return result
}
