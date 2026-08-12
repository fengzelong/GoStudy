package auth

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"time"

	"GoStudy/internal/domain"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
)

type Manager struct {
	secret []byte
	ttl    time.Duration
}

// Claims 是 JWT 的业务声明，保存用户身份、角色与有效期。
type Claims struct {
	Subject   int64           `json:"sub"`
	Role      domain.UserRole `json:"role"`
	ExpiresAt int64           `json:"exp"`
	IssuedAt  int64           `json:"iat"`
}

// NewManager 创建 JWT 管理器，空配置会回退到便于本地运行的默认值。
func NewManager(secret string, ttl time.Duration) *Manager {
	if secret == "" {
		secret = "gostudy-dev-secret"
	}
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl}
}

// Generate 为指定用户签发 HS256 JWT；省略角色时使用普通用户角色。
func (m *Manager) Generate(userID int64, roles ...domain.UserRole) (string, error) {
	role := domain.UserRoleUser
	if len(roles) > 0 && roles[0] != "" {
		role = roles[0]
	}
	now := time.Now()
	claims := Claims{Subject: userID, Role: role, IssuedAt: now.Unix(), ExpiresAt: now.Add(m.ttl).Unix()}
	header, err := json.Marshal(map[string]string{"alg": "HS256", "typ": "JWT"})
	if err != nil {
		return "", err
	}
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}
	unsigned := base64.RawURLEncoding.EncodeToString(header) + "." + base64.RawURLEncoding.EncodeToString(payload)
	return unsigned + "." + m.sign(unsigned), nil
}

// Parse 校验 JWT 结构、HS256 签名和有效期，并返回业务声明。
func (m *Manager) Parse(token string) (Claims, error) {
	parts := strings.Split(token, ".")
	if len(parts) != 3 || !hmac.Equal([]byte(m.sign(parts[0]+"."+parts[1])), []byte(parts[2])) {
		return Claims{}, ErrInvalidToken
	}
	headerBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var header struct {
		Algorithm string `json:"alg"`
	}
	if err := json.Unmarshal(headerBytes, &header); err != nil || header.Algorithm != "HS256" {
		return Claims{}, ErrInvalidToken
	}
	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		return Claims{}, ErrInvalidToken
	}
	var claims Claims
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Subject <= 0 || claims.ExpiresAt < time.Now().Unix() {
		if err == nil && claims.ExpiresAt < time.Now().Unix() {
			return Claims{}, ErrExpiredToken
		}
		return Claims{}, ErrInvalidToken
	}
	if claims.Role == "" {
		claims.Role = domain.UserRoleUser
	}
	return claims, nil
}

func (m *Manager) sign(body string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(body))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
