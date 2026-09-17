package auth

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"strings"
	"sync"
	"time"

	"GoStudy/internal/domain"
)

var (
	ErrInvalidToken = errors.New("invalid token")
	ErrExpiredToken = errors.New("expired token")
	ErrRevokedToken = errors.New("revoked token")
)

type Manager struct {
	secret  []byte
	ttl     time.Duration
	mu      sync.Mutex
	revoked map[string]int64
}

// Claims 是 JWT 的业务声明，保存用户身份、角色与有效期。
type Claims struct {
	Subject   int64           `json:"sub"`
	Role      domain.UserRole `json:"role"`
	ExpiresAt int64           `json:"exp"`
	IssuedAt  int64           `json:"iat"`
	TokenID   string          `json:"jti"`
}

// NewManager 创建 JWT 管理器，空配置会回退到便于本地运行的默认值。
func NewManager(secret string, ttl time.Duration) *Manager {
	if secret == "" {
		secret = "gostudy-dev-secret"
	}
	if ttl <= 0 {
		ttl = 2 * time.Hour
	}
	return &Manager{secret: []byte(secret), ttl: ttl, revoked: make(map[string]int64)}
}

// Generate 为指定用户签发 HS256 JWT；省略角色时使用普通用户角色。
func (m *Manager) Generate(userID int64, roles ...domain.UserRole) (string, error) {
	role := domain.UserRoleUser
	if len(roles) > 0 && roles[0] != "" {
		role = roles[0]
	}
	now := time.Now()
	tokenID, err := newTokenID()
	if err != nil {
		return "", err
	}
	claims := Claims{Subject: userID, Role: role, IssuedAt: now.Unix(), ExpiresAt: now.Add(m.ttl).Unix(), TokenID: tokenID}
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
	if err := json.Unmarshal(payload, &claims); err != nil || claims.Subject <= 0 || claims.TokenID == "" || claims.ExpiresAt < time.Now().Unix() {
		if err == nil && claims.ExpiresAt < time.Now().Unix() {
			return Claims{}, ErrExpiredToken
		}
		return Claims{}, ErrInvalidToken
	}
	if claims.Role == "" {
		claims.Role = domain.UserRoleUser
	}
	if m.isRevoked(claims) {
		return Claims{}, ErrRevokedToken
	}
	return claims, nil
}

// Refresh 吊销当前有效 Token 并为同一用户签发新 Token，实现轮换续期。
func (m *Manager) Refresh(claims Claims) (string, error) {
	m.Revoke(claims)
	return m.Generate(claims.Subject, claims.Role)
}

// Revoke 记录 Token 的唯一标识，直到其自然过期；内存记录会在服务重启后清空。
func (m *Manager) Revoke(claims Claims) {
	if claims.TokenID == "" {
		return
	}
	m.mu.Lock()
	defer m.mu.Unlock()
	m.revoked[claims.TokenID] = claims.ExpiresAt
}

func (m *Manager) isRevoked(claims Claims) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	for tokenID, expiresAt := range m.revoked {
		if expiresAt < time.Now().Unix() {
			delete(m.revoked, tokenID)
		}
	}
	_, ok := m.revoked[claims.TokenID]
	return ok
}

func newTokenID() (string, error) {
	value := make([]byte, 16)
	if _, err := rand.Read(value); err != nil {
		return "", err
	}
	return base64.RawURLEncoding.EncodeToString(value), nil
}

func (m *Manager) sign(body string) string {
	mac := hmac.New(sha256.New, m.secret)
	mac.Write([]byte(body))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil))
}
