package auth

import (
	"errors"
	"time"

	shared_errors "github.com/Matheus-Lima-Moreira/financial-pocket/internal/shared/errors"
	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	accessTokenSecret    string
	refreshTokenSecret   string
	accessTokenDuration  time.Duration
	refreshTokenDuration time.Duration
}

func NewJWTManager(accessTokenSecret, refreshTokenSecret string) *JWTManager {
	return &JWTManager{
		accessTokenSecret:    accessTokenSecret,
		refreshTokenSecret:   refreshTokenSecret,
		accessTokenDuration:  time.Hour,
		refreshTokenDuration: time.Hour * 24 * 7,
	}
}

func (j *JWTManager) GenerateAccessToken(userID string) (string, *shared_errors.AppError) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "access",
		"exp":     time.Now().Add(j.accessTokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	accessToken, err := token.SignedString([]byte(j.accessTokenSecret))
	if err != nil {
		return "", shared_errors.NewBadRequest(err.Error())
	}
	return accessToken, nil
}

func (j *JWTManager) GenerateRefreshToken(userID string) (string, *shared_errors.AppError) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"type":    "refresh",
		"exp":     time.Now().Add(j.refreshTokenDuration).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	refreshToken, err := token.SignedString([]byte(j.refreshTokenSecret))
	if err != nil {
		return "", shared_errors.NewBadRequest(err.Error())
	}
	return refreshToken, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (jwt.MapClaims, *shared_errors.AppError) {
	return j.parseToken(tokenString, j.accessTokenSecret)
}

func (j *JWTManager) ValidateRefreshToken(tokenString string) (string, *shared_errors.AppError) {
	claims, err := j.parseToken(tokenString, j.refreshTokenSecret)
	if err != nil {
		return "", err
	}

	tokenType, ok := claims["type"].(string)
	if !ok || tokenType != "refresh" {
		return "", shared_errors.NewUnauthorized("invalid token")
	}

	userID, ok := claims["user_id"].(string)
	if !ok {
		return "", shared_errors.NewUnauthorized("invalid token")
	}

	return userID, nil
}

func (j *JWTManager) parseToken(tokenString string, secret string) (jwt.MapClaims, *shared_errors.AppError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, shared_errors.NewUnauthorized("invalid token")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, shared_errors.NewUnauthorized("expired token")
		}
		return nil, shared_errors.NewUnauthorized("invalid token")
	}

	if !token.Valid {
		return nil, shared_errors.NewUnauthorized("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, shared_errors.NewUnauthorized("invalid token")
	}

	return claims, nil
}
