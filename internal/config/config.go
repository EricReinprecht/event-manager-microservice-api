package config

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port string
	Env  string

	FrontendURL string

	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string

	JWTSecret            string
	RefreshTokenDuration time.Duration

	TicketVerificationTTL time.Duration

	PayPalClientID     string
	PayPalClientSecret string
	PayPalBaseURL      string
	PayPalWebhookID    string

	PayPalReturnURL string
	PayPalCancelURL string

	CORSAllowedOrigins []string

	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string
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

		FrontendURL: getEnv(
			"FRONTEND_URL",
			"http://localhost:5173",
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

		RefreshTokenDuration: getEnvDuration(
			"REFRESH_TOKEN_DURATION",
			30*24*time.Hour,
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

		PayPalWebhookID: getEnv(
			"PAYPAL_WEBHOOK_ID",
			"",
		),

		PayPalReturnURL: getEnv(
			"PAYPAL_RETURN_URL",
			"",
		),

		PayPalCancelURL: getEnv(
			"PAYPAL_CANCEL_URL",
			"",
		),

		CORSAllowedOrigins: getEnvList(
			"CORS_ALLOWED_ORIGINS",
			[]string{
				"http://localhost:5173",
			},
		),

		SMTPHost: getEnv(
			"SMTP_HOST",
			"localhost",
		),

		SMTPPort: getEnvInt(
			"SMTP_PORT",
			587,
		),

		SMTPUser: getEnv(
			"SMTP_USER",
			"",
		),

		SMTPPassword: getEnv(
			"SMTP_PASSWORD",
			"",
		),

		SMTPFrom: getEnv(
			"SMTP_FROM",
			"",
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

func getEnvList(
	key string,
	fallback []string,
) []string {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	parts := strings.Split(
		value,
		",",
	)

	result := make([]string, 0, len(parts))

	for _, part := range parts {

		origin := strings.TrimSpace(part)

		if origin != "" {
			result = append(
				result,
				origin,
			)
		}
	}

	if len(result) == 0 {
		return fallback
	}

	return result
}

func getEnvDuration(
	key string,
	fallback time.Duration,
) time.Duration {

	value := os.Getenv(key)

	if value == "" {
		return fallback
	}

	duration, err := time.ParseDuration(value)

	if err != nil {
		return fallback
	}

	return duration
}
