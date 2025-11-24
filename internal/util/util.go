package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var AccessSecret = []byte("hotel_access_secret_key_2025_jwt")
var RefreshSecret = []byte("hotel_refresh_secret_key_2025_jwt")

type AccessClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

type RefreshClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	jwt.RegisteredClaims
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userId uint, userName string) (string, error) {
	claims := AccessClaims{
		UserId:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Minute * 15)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(AccessSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userId uint, userName string) (string, error) {
	claims := RefreshClaims{
		UserId:   userId,
		UserName: userName,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(RefreshSecret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// GenerateTokenPair 生成访问令牌和刷新令牌
func GenerateTokenPair(userId uint, userName string) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userId, userName)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(userId, userName)
	if err != nil {
		return "", "", err
	}
	return
}

// ParseAccessToken 解析访问令牌
func ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return AccessSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*AccessClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil

}

// ParseRefreshToken 解析刷新令牌
func ParseRefreshToken(tokenString string) (*RefreshClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &RefreshClaims{}, func(token *jwt.Token) (interface{}, error) {
		return RefreshSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil

}
