package utils

import (
	"crypto/rand"
	"encoding/hex"
	"time"

	"github.com/google/uuid"
)

// GenerateUUID 生成 UUID
func GenerateUUID() string {
	return uuid.New().String()
}

// GenerateToken 生成 Token (32字节随机数，hex编码)
func GenerateToken() (string, error) {
	b := make([]byte, 32)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// GenerateSalt 生成盐值
func GenerateSalt() ([]byte, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

// CalculateTokenExpiry 计算Token过期时间
func CalculateTokenExpiry(duration time.Duration) time.Time {
	return time.Now().Add(duration)
}
