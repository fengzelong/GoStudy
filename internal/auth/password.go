package auth

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
)

var ErrInvalidPasswordHash = errors.New("invalid password hash")

// HashPassword 使用标准库生成演示用密码摘要。
// 真实生产项目建议替换为 bcrypt、argon2 或统一认证中心。
func HashPassword(password string) (string, error) {
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	sum := sha256.Sum256(append(salt, []byte(password)...))
	return fmt.Sprintf("sha256$%s$%s", hex.EncodeToString(salt), hex.EncodeToString(sum[:])), nil
}

// VerifyPassword 使用恒定时间比较校验密码，避免把摘要字段暴露给接口层。
func VerifyPassword(encoded string, password string) (bool, error) {
	parts := strings.Split(encoded, "$")
	if len(parts) != 3 || parts[0] != "sha256" {
		return false, ErrInvalidPasswordHash
	}

	salt, err := hex.DecodeString(parts[1])
	if err != nil {
		return false, err
	}

	expected, err := hex.DecodeString(parts[2])
	if err != nil {
		return false, err
	}

	sum := sha256.Sum256(append(salt, []byte(password)...))
	return subtle.ConstantTimeCompare(sum[:], expected) == 1, nil
}
