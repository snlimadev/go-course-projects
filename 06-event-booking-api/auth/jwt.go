package auth

import (
	"errors"
	"time"

	"example.com/event-booking-api/config"
	"github.com/golang-jwt/jwt/v5"
)

type RefreshClaims struct {
	UserID int64 `json:"sub"`
	jwt.RegisteredClaims
}

type AccessClaims struct {
	UserID int64  `json:"sub"`
	Name   string `json:"name"`
	Email  string `json:"email"`
	jwt.RegisteredClaims
}

func GenerateRefreshToken(userID int64) (string, error) {
	duration := time.Duration(config.JWTRefreshDurationDays) * 24 * time.Hour

	claims := jwt.MapClaims{
		"sub": userID,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(duration).Unix(),
	}

	return signToken(claims, config.JWTRefreshSecret)
}

func GenerateAccessToken(userID int64, name, email string) (string, error) {
	duration := time.Duration(config.JWTAccessDurationMinutes) * time.Minute

	claims := jwt.MapClaims{
		"sub":   userID,
		"name":  name,
		"email": email,
		"iat":   time.Now().Unix(),
		"exp":   time.Now().Add(duration).Unix(),
	}

	return signToken(claims, config.JWTAccessSecret)
}

func ValidateRefreshToken(tokenString string) (*RefreshClaims, error) {
	claims, err := parseToken(tokenString, config.JWTRefreshSecret, &RefreshClaims{})

	if err != nil {
		return nil, err
	}

	return claims.(*RefreshClaims), nil
}

func ValidateAccessToken(tokenString string) (*AccessClaims, error) {
	claims, err := parseToken(tokenString, config.JWTAccessSecret, &AccessClaims{})

	if err != nil {
		return nil, err
	}

	return claims.(*AccessClaims), nil
}

func signToken(claims jwt.Claims, secret string) (string, error) {
	if secret == "" {
		return "", errors.New("JWT secret is not set for signing")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

func parseToken(tokenString, secret string, claims jwt.Claims) (jwt.Claims, error) {
	if secret == "" {
		return nil, errors.New("JWT secret is not set for parsing")
	}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}

		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)

	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, errors.New("invalid token")
	}

	return token.Claims, nil
}
