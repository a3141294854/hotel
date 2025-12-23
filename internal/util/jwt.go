package util

import (
	"time"

	"github.com/golang-jwt/jwt/v5"
	"hotel/internal/config"
)

var Secret []byte

var AccessExpireTime time.Duration
var RefreshExpireTime time.Duration

// AccessClaims 访问令牌声明
type AccessClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	HotelId  uint   `json:"hotel_id"`
	jwt.RegisteredClaims
}

// RefreshClaims 刷新令牌声明
type RefreshClaims struct {
	UserId   uint   `json:"user_id"`
	UserName string `json:"user_name"`
	HotelID  uint   `json:"hotel_id"`
	jwt.RegisteredClaims
}

// ConfigJwt 配置JWT
func ConfigJwt(cfg *config.Config) {
	Secret = []byte(cfg.JWT.SecretKey)
	AccessExpireTime = cfg.JWT.AccessTokenDuration
	RefreshExpireTime = cfg.JWT.RefreshTokenDuration
}

// GenerateAccessToken 生成访问令牌
func GenerateAccessToken(userId uint, userName string, HotelID uint) (string, error) {
	claims := AccessClaims{
		UserId:   userId,
		UserName: userName,
		HotelId:  HotelID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(AccessExpireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// GenerateRefreshToken 生成刷新令牌
func GenerateRefreshToken(userId uint, userName string, HotelID uint) (string, error) {
	claims := RefreshClaims{
		UserId:   userId,
		UserName: userName,
		HotelID:  HotelID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(RefreshExpireTime)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(Secret)
	if err != nil {
		return "", err
	}
	return tokenString, nil

}

// GenerateTokenPair 生成访问令牌和刷新令牌
func GenerateTokenPair(userId uint, userName string, HotelID uint) (accessToken string, refreshToken string, err error) {
	accessToken, err = GenerateAccessToken(userId, userName, HotelID)
	if err != nil {
		return "", "", err
	}
	refreshToken, err = GenerateRefreshToken(userId, userName, HotelID)
	if err != nil {
		return "", "", err
	}
	return
}

// ParseAccessToken 解析访问令牌
func ParseAccessToken(tokenString string) (*AccessClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &AccessClaims{}, func(token *jwt.Token) (interface{}, error) {
		return Secret, nil
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
		return Secret, nil
	})
	if err != nil {
		return nil, err
	}
	if claims, ok := token.Claims.(*RefreshClaims); ok && token.Valid {
		return claims, nil
	}
	return nil, nil

}
