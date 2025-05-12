package security

import (
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"strings"
	"task-api/pkg/config"
	"time"
)

type JWTClaims struct {
	UserID uuid.UUID `json:"user_id"`
	jwt.RegisteredClaims
}

func CreateAccessJWT(cfg config.AppConfig, userID uuid.UUID) (string, error) {
	claims := JWTClaims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(cfg.Auth.JWTExpiry)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod(cfg.Auth.JWTAlgorithm), claims)
	signedToken, err := token.SignedString([]byte(cfg.Auth.JWTSecret))
	if err != nil {
		return "", err
	}
	return signedToken, nil
}

func ParseAccessJWT(cfg config.AppConfig, tokenStr string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenStr, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		if token.Method.Alg() != cfg.Auth.JWTAlgorithm {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Auth.JWTSecret), nil
	})

	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, err
	}
	claims := token.Claims.(*JWTClaims)
	return claims, nil
}

func TokenString(auth string) (string, error) {
	parts := strings.SplitN(auth, " ", 2)
	if !(len(parts) == 2 && parts[0] == "Bearer") {
		return "", errors.New("unauthorized")
	}
	tokenStr := parts[1]
	return tokenStr, nil
}
