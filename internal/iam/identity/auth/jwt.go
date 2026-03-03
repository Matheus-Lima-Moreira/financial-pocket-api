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

func (j *JWTManager) GenerateAccessToken(userID, organizationID string) (string, *shared_errors.AppError) {
	claims := NewJWTManagerClaims(userID, organizationID, "access", j.accessTokenDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.ToJwtMapClaims())
	accessToken, err := token.SignedString([]byte(j.accessTokenSecret))
	if err != nil {
		return "", shared_errors.NewBadRequest(err.Error())
	}
	return accessToken, nil
}

func (j *JWTManager) GenerateRefreshToken(userID, organizationID string) (string, *shared_errors.AppError) {
	claims := NewJWTManagerClaims(userID, organizationID, "refresh", j.refreshTokenDuration)
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims.ToJwtMapClaims())
	refreshToken, err := token.SignedString([]byte(j.refreshTokenSecret))
	if err != nil {
		return "", shared_errors.NewBadRequest(err.Error())
	}
	return refreshToken, nil
}

func (j *JWTManager) ValidateToken(tokenString string) (jwt.MapClaims, *shared_errors.AppError) {
	return j.parseToken(tokenString, j.accessTokenSecret)
}

func (j *JWTManager) ValidateRefreshToken(tokenString string) (*JWTManagerClaims, *shared_errors.AppError) {
	claims, err := j.parseToken(tokenString, j.refreshTokenSecret)
	if err != nil {
		return nil, err
	}

	claimsObj, err := ParseJWTManagerClaims(claims)
	if err != nil {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}
	if claimsObj.TokenType != "refresh" {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}

	return claimsObj, nil
}

func (j *JWTManager) parseToken(tokenString string, secret string) (jwt.MapClaims, *shared_errors.AppError) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, shared_errors.NewUnauthorized("error.invalid_token")
		}
		return []byte(secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, shared_errors.NewUnauthorized("error.expired_token")
		}
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}

	if !token.Valid {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}

	return claims, nil
}

type JWTManagerClaims struct {
	UserID         string `json:"user_id"`
	OrganizationID string `json:"organization_id"`
	TokenType      string `json:"type"`
	Exp            int64  `json:"exp"`
	Iat            int64  `json:"iat"`
}

func NewJWTManagerClaims(userID, organizationID, tokenType string, accessTokenDuration time.Duration) JWTManagerClaims {
	return JWTManagerClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		TokenType:      tokenType,
		Exp:            time.Now().Add(accessTokenDuration).Unix(),
		Iat:            time.Now().Unix(),
	}
}

func (c *JWTManagerClaims) ToJwtMapClaims() jwt.MapClaims {
	return jwt.MapClaims{
		"user_id":         c.UserID,
		"organization_id": c.OrganizationID,
		"type":            c.TokenType,
		"exp":             c.Exp,
		"iat":             c.Iat,
	}
}

func ParseJWTManagerClaims(claims jwt.MapClaims) (*JWTManagerClaims, *shared_errors.AppError) {
	userID, ok := claims["user_id"].(string)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}
	organizationID, ok := claims["organization_id"].(string)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}
	tokenType, ok := claims["type"].(string)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}
	exp, ok := claims["exp"].(int64)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}
	iat, ok := claims["iat"].(int64)
	if !ok {
		return nil, shared_errors.NewUnauthorized("error.invalid_token")
	}

	return &JWTManagerClaims{
		UserID:         userID,
		OrganizationID: organizationID,
		TokenType:      tokenType,
		Exp:            exp,
		Iat:            iat,
	}, nil
}
