package jwt

import (
	"crypto/rsa"
	"errors"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type Config struct {
	Algorithm      string
	Secret         string
	PrivateKeyPath string
	PublicKeyPath  string
	ExpireHours    int
}

var cfg *Config
var rsaPrivateKey *rsa.PrivateKey
var rsaPublicKey *rsa.PublicKey

type Claims struct {
	UserID   string `json:"user_id"`
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func Init(c *Config) error {
	cfg = c
	if cfg == nil {
		return errors.New("jwt config is nil")
	}

	if cfg.Algorithm == "RS256" {
		// 读取私钥
		privBytes, err := os.ReadFile(cfg.PrivateKeyPath)
		if err != nil {
			return err
		}
		privKey, err := jwt.ParseRSAPrivateKeyFromPEM(privBytes)
		if err != nil {
			return err
		}
		rsaPrivateKey = privKey

		// 读取公钥
		pubBytes, err := os.ReadFile(cfg.PublicKeyPath)
		if err != nil {
			return err
		}
		pubKey, err := jwt.ParseRSAPublicKeyFromPEM(pubBytes)
		if err != nil {
			return err
		}
		rsaPublicKey = pubKey
	}
	return nil
}

// GenerateToken 生成 Token
func GenerateToken(userID, username string) (string, error) {
	if cfg == nil {
		return "", errors.New("jwt not initialized")
	}
	exp := time.Hour * time.Duration(cfg.ExpireHours)
	now := time.Now()

	claims := Claims{
		UserID:   userID,
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(exp)), // 24小时过期
			Issuer:    "blog_go",                               // 签发者
		},
	}

	if cfg.Algorithm == "RS256" {
		token := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
		return token.SignedString(rsaPrivateKey)
	}

	// 默认 HS256
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(cfg.Secret))
}

// ParseToken 解析并验证Token
func ParseToken(tokenString string) (*Claims, error) {
	if cfg == nil {
		return nil, errors.New("jwt not initialized")
	}
	claims := &Claims{}

	keyFunc := func(token *jwt.Token) (interface{}, error) {
		// 根据算法验证
		if cfg.Algorithm == "RS256" {
			if _, ok := token.Method.(*jwt.SigningMethodRSA); !ok {
				return nil, errors.New("unexpected signing method")
			}
			return rsaPublicKey, nil
		}
		// HS
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(cfg.Secret), nil
	}

	parsed, err := jwt.ParseWithClaims(tokenString, claims, keyFunc)
	if err != nil {
		return nil, err
	}
	if !parsed.Valid {
		return nil, errors.New("invalid token")
	}
	return claims, nil
}
