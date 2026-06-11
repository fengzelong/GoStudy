package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"GoStudy/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

// userModel 是 GORM 持久化模型，和 domain.User 分开以隔离数据库细节。
type userModel struct {
	ID           int64     `gorm:"primaryKey"`
	Name         string    `gorm:"size:100;not null"`
	Email        string    `gorm:"size:255;uniqueIndex;not null"`
	PasswordHash string    `gorm:"size:255;not null"`
	CreatedAt    time.Time `gorm:"not null"`
}

// taskModel 是任务表结构，领域层只依赖 domain.Task。
type taskModel struct {
	ID        int64             `gorm:"primaryKey"`
	Title     string            `gorm:"size:255;not null"`
	OwnerID   int64             `gorm:"index;not null"`
	Status    domain.TaskStatus `gorm:"size:32;not null"`
	CreatedAt time.Time         `gorm:"not null"`
	UpdatedAt time.Time         `gorm:"not null"`
}

// NewGormStore 建立 MySQL 连接并自动迁移企业骨架所需表结构。
func NewGormStore(dsn string) (*GormStore, error) {
	if dsn == "" {
		return nil, fmt.Errorf("mysql dsn is required")
	}

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	if err := db.AutoMigrate(&userModel{}, &taskModel{}); err != nil {
		return nil, err
	}

	return &GormStore{db: db}, nil
}

// Close 释放 GORM 底层数据库连接池。
func (s *GormStore) Close() error {
	sqlDB, err := s.db.DB()
	if err != nil {
		return err
	}
	return sqlDB.Close()
}

// CreateUser 持久化用户并返回带数据库主键的领域对象。
func (s *GormStore) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	model := userModel{
		Name:         user.Name,
		Email:        user.Email,
		PasswordHash: user.PasswordHash,
		CreatedAt:    time.Now(),
	}

	if err := s.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.User{}, err
	}
	return model.toDomain(), nil
}

// GetUser 按主键查询用户。
func (s *GormStore) GetUser(ctx context.Context, id int64) (domain.User, error) {
	var model userModel
	if err := s.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.User{}, mapGormError(err)
	}
	return model.toDomain(), nil
}

// FindUserByEmail 支持登录和注册去重场景。
func (s *GormStore) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	var model userModel
	if err := s.db.WithContext(ctx).Where("email = ?", email).First(&model).Error; err != nil {
		return domain.User{}, mapGormError(err)
	}
	return model.toDomain(), nil
}

// ListUsers 返回用户列表，示例中按 ID 保持稳定顺序。
func (s *GormStore) ListUsers(ctx context.Context) ([]domain.User, error) {
	var models []userModel
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	users := make([]domain.User, 0, len(models))
	for _, model := range models {
		users = append(users, model.toDomain())
	}
	return users, nil
}

// CreateTask 创建任务并补齐创建、更新时间。
func (s *GormStore) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	now := time.Now()
	model := taskModel{
		Title:     task.Title,
		OwnerID:   task.OwnerID,
		Status:    task.Status,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.db.WithContext(ctx).Create(&model).Error; err != nil {
		return domain.Task{}, err
	}
	return model.toDomain(), nil
}

// GetTask 按主键查询任务。
func (s *GormStore) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	var model taskModel
	if err := s.db.WithContext(ctx).First(&model, id).Error; err != nil {
		return domain.Task{}, mapGormError(err)
	}
	return model.toDomain(), nil
}

// ListTasks 返回任务列表，示例中按 ID 保持稳定顺序。
func (s *GormStore) ListTasks(ctx context.Context) ([]domain.Task, error) {
	var models []taskModel
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	tasks := make([]domain.Task, 0, len(models))
	for _, model := range models {
		tasks = append(tasks, model.toDomain())
	}
	return tasks, nil
}

// UpdateTask 只更新已存在任务，避免 Save 意外插入新记录。
func (s *GormStore) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	var model taskModel
	if err := s.db.WithContext(ctx).First(&model, task.ID).Error; err != nil {
		return domain.Task{}, mapGormError(err)
	}

	model.Title = task.Title
	model.OwnerID = task.OwnerID
	model.Status = task.Status
	model.UpdatedAt = time.Now()

	if err := s.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.Task{}, err
	}
	return model.toDomain(), nil
}

func (m userModel) toDomain() domain.User {
	return domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		PasswordHash: m.PasswordHash,
		CreatedAt:    m.CreatedAt,
	}
}

func (m taskModel) toDomain() domain.Task {
	return domain.Task{
		ID:        m.ID,
		Title:     m.Title,
		OwnerID:   m.OwnerID,
		Status:    m.Status,
		CreatedAt: m.CreatedAt,
		UpdatedAt: m.UpdatedAt,
	}
}

func mapGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
