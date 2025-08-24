package utils

import (
	"fmt"
	"time"

	"go-crud-employee/config"
	"go-crud-employee/models"

	"github.com/golang-jwt/jwt/v5"
)

type JWTManager struct {
	secret string
	expiry time.Duration
}

func NewJWTManager(cfg *config.Config) *JWTManager {
	return &JWTManager{
		secret: cfg.JWT.Secret,
		expiry: cfg.JWT.Expiry,
	}
}

// GenerateToken generates a new JWT token for the user
func (j *JWTManager) GenerateToken(user *models.User) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(j.expiry)

	claims := models.JWTClaims{
		UserID:   user.ID,
		Username: user.Username,
		Email:    user.Email,
		Exp:      expiresAt.Unix(),
		Iat:      now.Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  claims.UserID,
		"username": claims.Username,
		"email":    claims.Email,
		"exp":      claims.Exp,
		"iat":      claims.Iat,
	})

	tokenString, err := token.SignedString([]byte(j.secret))
	if err != nil {
		return "", time.Time{}, fmt.Errorf("failed to sign token: %v", err)
	}

	return tokenString, expiresAt, nil
}

// ValidateToken validates and parses a JWT token
func (j *JWTManager) ValidateToken(tokenString string) (*models.JWTClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(j.secret), nil
	})

	if err != nil {
		return nil, fmt.Errorf("failed to parse token: %v", err)
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token claims")
	}

	// Check expiration
	exp, ok := claims["exp"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid expiration claim")
	}

	if time.Now().Unix() > int64(exp) {
		return nil, fmt.Errorf("token has expired")
	}

	// Extract claims
	userID, ok := claims["user_id"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid user_id claim")
	}

	username, ok := claims["username"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid username claim")
	}

	email, ok := claims["email"].(string)
	if !ok {
		return nil, fmt.Errorf("invalid email claim")
	}

	iat, ok := claims["iat"].(float64)
	if !ok {
		return nil, fmt.Errorf("invalid iat claim")
	}

	return &models.JWTClaims{
		UserID:   int(userID),
		Username: username,
		Email:    email,
		Exp:      int64(exp),
		Iat:      int64(iat),
	}, nil
}