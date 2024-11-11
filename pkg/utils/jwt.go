package utils

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v4"
	"github.com/spf13/viper"
)

// TokenType 定义token类型
type TokenType string

const (
	AccessToken  TokenType = "access"
	RefreshToken TokenType = "refresh"
)

// CustomClaims 自定义JWT claims
type CustomClaims struct {
	UserID    uint      `json:"user_id"`
	Username  string    `json:"username"`
	TokenType TokenType `json:"token_type"`
	jwt.RegisteredClaims
}

// GenerateToken 生成JWT token
func GenerateToken(userID uint, username string, tokenType TokenType) (string, error) {
	var secret string
	var expire time.Duration

	if tokenType == AccessToken {
		secret = viper.GetString("jwt.access_secret")
		expire = time.Duration(viper.GetInt("jwt.access_expire")) * time.Second
	} else {
		secret = viper.GetString("jwt.refresh_secret")
		expire = time.Duration(viper.GetInt("jwt.refresh_expire")) * time.Second
	}

	claims := CustomClaims{
		UserID:    userID,
		Username:  username,
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// GenerateTokenWithClaims 生成带自定义claims的JWT token
func GenerateTokenWithClaims(claims map[string]interface{}, tokenType TokenType) (string, error) {
	var secret string
	var expire time.Duration

	if tokenType == AccessToken {
		secret = viper.GetString("jwt.access_secret")
		expire = time.Duration(viper.GetInt("jwt.access_expire")) * time.Second
	} else {
		secret = viper.GetString("jwt.refresh_secret")
		expire = time.Duration(viper.GetInt("jwt.refresh_expire")) * time.Second
	}

	// 创建自定义claims
	customClaims := CustomClaims{
		UserID:    claims["user_id"].(uint),
		Username:  claims["username"].(string),
		TokenType: tokenType,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(expire)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, customClaims)
	return token.SignedString([]byte(secret))
}

// ParseToken 解析JWT token
func ParseToken(tokenString string, tokenType TokenType) (*CustomClaims, error) {
	var secret string
	if tokenType == AccessToken {
		secret = viper.GetString("jwt.access_secret")
	} else {
		secret = viper.GetString("jwt.refresh_secret")
	}

	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*CustomClaims); ok && token.Valid {
		if claims.TokenType != tokenType {
			return nil, errors.New("invalid token type")
		}
		return claims, nil
	}

	return nil, errors.New("invalid token")
}
