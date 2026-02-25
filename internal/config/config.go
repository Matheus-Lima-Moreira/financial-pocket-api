package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Env                string
	Port               string
	DBUser             string
	DBPassword         string
	DBHost             string
	DBPort             string
	DBName             string
	AccessTokenSecret  string
	RefreshTokenSecret string
	SMTPHost           string
	SMTPPort           string
	SMTPUser           string
	SMTPPassword       string
	SMTPFrom           string
	VerifyEmailBaseURL string
}

func Load() *Config {
	if os.Getenv("ENV") != "production" {
		_ = godotenv.Load()
	}

	env := os.Getenv("ENV")
	if env == "" {
		env = "development"
	}

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	dbUser := os.Getenv("DB_USER")
	if dbUser == "" {
		dbUser = "root"
	}

	dbPassword, ok := os.LookupEnv("DB_PASSWORD")
	if !ok {
		dbPassword = "root"
	}

	dbHost := os.Getenv("DB_HOST")
	if dbHost == "" {
		dbHost = "localhost"
	}

	dbPort := os.Getenv("DB_PORT")
	if dbPort == "" {
		dbPort = "3306"
	}

	dbName := os.Getenv("DB_NAME")
	if dbName == "" {
		dbName = "financial_pocket"
	}

	accessTokenSecret := os.Getenv("ACCESS_TOKEN_SECRET")
	if accessTokenSecret == "" {
		accessTokenSecret = "secret"
	}

	refreshTokenSecret := os.Getenv("REFRESH_TOKEN_SECRET")
	if refreshTokenSecret == "" {
		refreshTokenSecret = "secret"
	}

	smtpHost := os.Getenv("SMTP_HOST")
	if smtpHost == "" {
		smtpHost = "localhost"
	}

	smtpPort := os.Getenv("SMTP_PORT")
	if smtpPort == "" {
		smtpPort = "1025"
	}

	smtpUser := os.Getenv("SMTP_USER")
	smtpPassword := os.Getenv("SMTP_PASSWORD")

	smtpFrom := os.Getenv("SMTP_FROM")
	if smtpFrom == "" {
		smtpFrom = "noreply@financial-pocket.dev"
	}

	verifyEmailBaseURL := os.Getenv("VERIFY_EMAIL_BASE_URL")
	if verifyEmailBaseURL == "" {
		verifyEmailBaseURL = "http://localhost:8090/auth/verify-email"
	}

	return &Config{
		Env:                env,
		Port:               port,
		DBUser:             dbUser,
		DBPassword:         dbPassword,
		DBHost:             dbHost,
		DBPort:             dbPort,
		DBName:             dbName,
		AccessTokenSecret:  accessTokenSecret,
		RefreshTokenSecret: refreshTokenSecret,
		SMTPHost:           smtpHost,
		SMTPPort:           smtpPort,
		SMTPUser:           smtpUser,
		SMTPPassword:       smtpPassword,
		SMTPFrom:           smtpFrom,
		VerifyEmailBaseURL: verifyEmailBaseURL,
	}
}

func (c *Config) TrustedProxies() []string {
	if c.Env == "development" {
		return []string{"127.0.0.1", "::1"}
	}

	trustedProxies := os.Getenv("TRUSTED_PROXIES")
	if trustedProxies == "" {
		return nil
	}

	proxies := strings.Split(trustedProxies, ",")
	result := make([]string, 0, len(proxies))
	for _, proxy := range proxies {
		trimmed := strings.TrimSpace(proxy)
		if trimmed != "" {
			result = append(result, trimmed)
		}
	}

	if len(result) == 0 {
		return nil
	}

	return result
}

func (c *Config) DSN() string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		c.DBUser,
		c.DBPassword,
		c.DBHost,
		c.DBPort,
		c.DBName,
	)
}
