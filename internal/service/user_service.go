package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"GoStudy/internal/auth"
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
}

// NewUserService 创建用户服务，业务层只依赖仓储接口。
func NewUserService(users repository.UserRepository) *UserService {
	return &UserService{users: users}
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

// List 返回用户列表，后续可以在这里加入分页和权限规则。
func (s *UserService) List(ctx context.Context) ([]domain.User, error) {
	return s.users.ListUsers(ctx)
}
