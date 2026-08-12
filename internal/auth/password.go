package auth

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

var ErrInvalidPasswordHash = errors.New("invalid password hash")

// HashPassword 使用 bcrypt 生成密码摘要，避免快速散列被暴力破解。
func HashPassword(password string) (string, error) {
	encoded, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(encoded), err
}

// VerifyPassword 使用 bcrypt 校验密码。
func VerifyPassword(encoded string, password string) (bool, error) {
	if err := bcrypt.CompareHashAndPassword([]byte(encoded), []byte(password)); err != nil {
		if errors.Is(err, bcrypt.ErrMismatchedHashAndPassword) {
			return false, nil
		}
		return false, ErrInvalidPasswordHash
	}
	return true, nil
}
