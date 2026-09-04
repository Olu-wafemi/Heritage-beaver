package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const (
	TokenIssuer      = "heritage-weaver"
	AccessTokenTTL   = 30 * time.Minute
	RefreshTokenTTL  = 30 * 24 * time.Hour
	MinPasswordChars = 8
)

type Claims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(userID string, secret string) (string, error) {
	now := time.Now()
	claims := Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    TokenIssuer,
			ExpiresAt: jwt.NewNumericDate(now.Add(AccessTokenTTL)),
			IssuedAt:  jwt.NewNumericDate(now),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := token.SignedString([]byte(secret))
	if err != nil {
		return "", fmt.Errorf("sign token: %w", err)
	}

	return signed, nil
}

// GenerateToken is kept for compatibility and delegates to GenerateAccessToken.
func GenerateToken(userID string, secret string) (string, error) {
	return GenerateAccessToken(userID, secret)
}

// GenerateRefreshToken returns an opaque token and its SHA-256 hash.
// Only the hash is stored; the token itself is shown once to the client.
func GenerateRefreshToken() (token string, hash string, err error) {
	raw := make([]byte, 32)
	if _, err := rand.Read(raw); err != nil {
		return "", "", fmt.Errorf("random refresh token: %w", err)
	}

	token = hex.EncodeToString(raw)
	return token, HashRefreshToken(token), nil
}

func HashRefreshToken(token string) string {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:])
}

func ValidateToken(tokenString string, secret string) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(secret), nil
	})
	if err != nil {
		return nil, fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}

	if claims.Issuer != "" && claims.Issuer != TokenIssuer {
		return nil, fmt.Errorf("unexpected token issuer")
	}

	if claims.UserID == "" {
		return nil, fmt.Errorf("token missing subject")
	}

	return claims, nil
}

func GetSecret() string {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return "dev-secret-do-not-use-in-production"
	}
	return secret
}
