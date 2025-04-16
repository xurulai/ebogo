package user

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"strings"
)

// HashPassword 使用 SHA-256 算法对密码进行加密
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := sha256.Sum256([]byte(password + string(salt)))

	// 将盐值和哈希值编码为 Base64 字符串
	saltB64 := base64.StdEncoding.EncodeToString(salt)
	hashB64 := base64.StdEncoding.EncodeToString(hash[:])

	return saltB64 + "$" + hashB64, nil
}

// CheckPassword 使用 SHA-256 算法验证密码
func CheckPassword(hashedPassword, password string) error {
	parts := strings.Split(hashedPassword, "$")
	if len(parts) != 2 {
		return errors.New("invalid hashed password format")
	}

	salt, err := base64.StdEncoding.DecodeString(parts[0])
	if err != nil {
		return err
	}

	hash := sha256.Sum256([]byte(password + string(salt)))
	hashB64 := base64.StdEncoding.EncodeToString(hash[:])

	if hashB64 != parts[1] {
		return errors.New("password mismatch")
	}

	return nil
}