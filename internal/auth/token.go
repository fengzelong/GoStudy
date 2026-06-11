package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

// Claims 是业务 Token 的最小声明，当前只保存用户 ID 和有效期。
type Claims struct {
	Subject   int64 `json:"sub"`
	ExpiresAt int64 `json:"exp"`
	IssuedAt  int64 `json:"iat"`
}

// NewManager 创建 Token 管理器，空配置会回退到便于本地运行的默认值。
func NewManager(secret string, ttl time.Duration) *Manager {
	if secret == "" {
		secret = "gostudy-dev-secret"
	}
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Generate 为指定用户签发带 HMAC 签名的轻量 Token。
func (m *Manager) Generate(userID int64) (string, error) {
	now := time.Now()
	claims := Claims{
		Subject:   userID,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(m.ttl).Unix(),
	}

	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	body := base64.RawURLEncoding.EncodeToString(payload)
	signature := m.sign(body)
	return body + "." + signature, nil
}

// Parse 校验 Token 签名和有效期，并返回业务声明。
func (m *Manager) Parse(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return Claims{}, ErrInvalidToken
	}

	if !hmac.Equal([]byte(m.sign(parts[0])), []byte(parts[1])) {
		return Claims{}, ErrInvalidToken
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}

	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil {
		return Claims{}, ErrInvalidToken
	}
	if claims.ExpiresAt < time.Now().Unix() {
		return Claims{}, ErrExpiredToken
	}
	return claims, nil
}

func (m *Manager) sign(body string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(body))
	mac.Write([]byte(strconv.FormatInt(int64(len(body)), 10)))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
