package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"time"

	"GoStudy/internal/auth"
	"GoStudy/internal/cache"
	"GoStudy/internal/domain"
	"GoStudy/internal/repository"
)

type RegisterUserInput struct {
	Name     string `json:"name" binding:"required"`
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginInput struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type UserService struct {
	users repository.UserRepository
	cache cache.Store
}

// NewUserService 创建用户服务，业务层只依赖仓储接口。
func NewUserService(users repository.UserRepository, stores ...cache.Store) *UserService {
	userCache := cache.NewNoop()
	if len(stores) > 0 && stores[0] != nil {
		userCache = stores[0]
	}
	return &UserService{users: users, cache: userCache}
}

// Register 注册用户，包含基础校验、邮箱去重和密码摘要。
func (s *UserService) Register(ctx context.Context, input RegisterUserInput) (domain.User, error) {
	name := strings.TrimSpace(input.Name)
	email := strings.ToLower(strings.TrimSpace(input.Email))
	if name == "" || email == "" || len(input.Password) < 6 {
		return domain.User{}, fmt.Errorf("%w: name, email and password are required", ErrInvalidInput)
	}

	if _, err := s.users.FindUserByEmail(ctx, email); err == nil {
		return domain.User{}, fmt.Errorf("%w: email already exists", ErrConflict)
	} else if !errors.Is(err, repository.ErrNotFound) {
		return domain.User{}, err
	}

	passwordHash, err := auth.HashPassword(input.Password)
	if err != nil {
		return domain.User{}, err
	}

	return s.users.CreateUser(ctx, domain.User{
		Name:         name,
		Email:        email,
		PasswordHash: passwordHash,
	})
}

// Login 校验邮箱和密码，成功后返回不包含密码摘要的 JSON 用户对象。
func (s *UserService) Login(ctx context.Context, input LoginInput) (domain.User, error) {
	email := strings.ToLower(strings.TrimSpace(input.Email))
	user, err := s.users.FindUserByEmail(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.User{}, fmt.Errorf("%w: email or password is wrong", ErrNotFound)
		}
		return domain.User{}, err
	}

	ok, err := auth.VerifyPassword(user.PasswordHash, input.Password)
	if err != nil {
		return domain.User{}, err
	}
	if !ok {
		return domain.User{}, fmt.Errorf("%w: email or password is wrong", ErrNotFound)
	}

	return user, nil
}

// UserPage 是用户列表及其分页信息。
type UserPage struct {
	PageMeta
	Items []domain.User `json:"items"`
}

// Get 返回指定用户，用于当前用户资料查询等场景。
func (s *UserService) Get(ctx context.Context, id int64) (domain.User, error) {
	key := fmt.Sprintf("user:%d", id)
	if value, ok, err := s.cache.Get(ctx, key); err == nil && ok {
		var user domain.User
		if err := json.Unmarshal(value, &user); err == nil {
			return user, nil
		}
	}
	user, err := s.users.GetUser(ctx, id)
	if errors.Is(err, repository.ErrNotFound) {
		return domain.User{}, fmt.Errorf("%w: user not found", ErrNotFound)
	}
	if err == nil {
		if value, marshalErr := json.Marshal(user); marshalErr == nil {
			_ = s.cache.Set(ctx, key, value, time.Minute)
		}
	}
	return user, err
}

// List 返回分页后的用户列表。
func (s *UserService) List(ctx context.Context, input PageInput) (UserPage, error) {
	users, err := s.users.ListUsers(ctx)
	if err != nil {
		return UserPage{}, err
	}
	meta, start, end, err := normalizePage(input, len(users))
	if err != nil {
		return UserPage{}, err
	}
	return UserPage{PageMeta: meta, Items: users[start:end]}, nil
}
