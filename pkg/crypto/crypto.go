package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"sync"
)

// RSAKeyPair 存储密钥对
type RSAKeyPair struct {
	PrivateKey *rsa.PrivateKey
	PublicKey  string // PEM 格式的公钥字符串
}

// 单例模式
var (
	GlobalKeyPair *RSAKeyPair
	once          sync.Once
)

// InitRSAKeyPair 初始化 RSA 密钥对 (2048位)
func InitRSAKeyPair() error {
	var err error
	once.Do(func() {
		RSAKeyPair, e := generateKeyPair(2048)
		GlobalKeyPair = RSAKeyPair
		err = e
	})
	return err
}

// generateKeyPair 生成私钥和公钥
func generateKeyPair(bits int) (*RSAKeyPair, error) {
	// 生成私钥
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, err
	}

	// 生成公钥 PEM 字符串
	publicKey := &privateKey.PublicKey
	pubASN1, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		return nil, err
	}

	pubBytes := pem.EncodeToMemory(&pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: pubASN1,
	})

	return &RSAKeyPair{
		PrivateKey: privateKey,
		PublicKey:  string(pubBytes),
	}, nil
}

// Decrypt 使用私钥解密
func Decrypt(encryptedBase64 string) (string, error) {
	if GlobalKeyPair == nil || GlobalKeyPair.PrivateKey == nil {
		return "", errors.New("RSA key pair not initialized")
	}

	// Base64 解码
	encryptedBytes, err := base64.StdEncoding.DecodeString(encryptedBase64)
	if err != nil {
		return "", errors.New("invalid base64 string")
	}

	// RSA 解密 (PKCS1v15 padding)
	decryptedBytes, err := rsa.DecryptPKCS1v15(rand.Reader, GlobalKeyPair.PrivateKey, encryptedBytes)
	if err != nil {
		return "", errors.New("decryption failed")
	}

	return string(decryptedBytes), nil
}

// GetPublicKey 获取公钥字符串
func GetPublicKey() (string, error) {
	if GlobalKeyPair == nil {
		return "", errors.New("RSA key pair not initialized")
	}
	return GlobalKeyPair.PublicKey, nil
}
