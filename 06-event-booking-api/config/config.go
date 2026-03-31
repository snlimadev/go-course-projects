package config

import (
	"os"
	"strconv"
)

var (
	Environment              string
	JWTAccessSecret          string
	JWTRefreshSecret         string
	JWTAccessDurationMinutes int
	JWTRefreshDurationDays   int
)

func Load() {
	Environment = getEnv("GO_ENV", "production")
	JWTAccessSecret = os.Getenv("JWT_ACCESS_SECRET")
	JWTRefreshSecret = os.Getenv("JWT_REFRESH_SECRET")
	JWTAccessDurationMinutes = getEnvInt("JWT_ACCESS_DURATION_MINUTES", 60)
	JWTRefreshDurationDays = getEnvInt("JWT_REFRESH_DURATION_DAYS", 30)
}

func getEnv(key, defaultValue string) string {
	value := os.Getenv(key)

	if value == "" {
		return defaultValue
	}

	return value
}

func getEnvInt(key string, defaultValue int) int {
	value, err := strconv.Atoi(os.Getenv(key))

	if err != nil || value <= 0 {
		return defaultValue
	}

	return value
}
