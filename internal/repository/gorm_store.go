package repository

import (
	"context"
	"errors"
	"fmt"
	"time"

	"GoStudy/internal/audit"
	"GoStudy/internal/domain"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type GormStore struct {
	db *gorm.DB
}

// userModel 是 GORM 持久化模型，和 domain.User 分开以隔离数据库细节。
type userModel struct {
	ID           int64           `gorm:"primaryKey"`
	Name         string          `gorm:"size:100;not null"`
	Email        string          `gorm:"size:255;uniqueIndex;not null"`
	Role         domain.UserRole `gorm:"size:32;not null"`
	PasswordHash string          `gorm:"size:255;not null"`
	CreatedAt    time.Time       `gorm:"not null"`
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

// auditModel 是审计记录的持久化模型，避免领域服务依赖 GORM。
type auditModel struct {
	ID        int64     `gorm:"primaryKey"`
	Action    string    `gorm:"size:100;not null"`
	ActorID   int64     `gorm:"index;not null"`
	Resource  string    `gorm:"size:255;not null"`
	CreatedAt time.Time `gorm:"index;not null"`
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

	if err := db.AutoMigrate(&userModel{}, &taskModel{}, &auditModel{}); err != nil {
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

// Health 检查底层 MySQL 连接池是否可用。
func (s *GormStore) Health(ctx context.Context) string {
	sqlDB, err := s.db.DB()
	if err != nil || sqlDB.PingContext(ctx) != nil {
		return "down"
	}
	return "ok"
}

// CreateUser 持久化用户并返回带数据库主键的领域对象。
func (s *GormStore) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	model := userModel{
		Name:         user.Name,
		Email:        user.Email,
		Role:         user.Role,
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

// UpdateUser 更新已有用户，用于角色调整等管理操作。
func (s *GormStore) UpdateUser(ctx context.Context, user domain.User) (domain.User, error) {
	var model userModel
	if err := s.db.WithContext(ctx).First(&model, user.ID).Error; err != nil {
		return domain.User{}, mapGormError(err)
	}
	model.Name = user.Name
	model.Email = user.Email
	model.Role = user.Role
	model.PasswordHash = user.PasswordHash
	if err := s.db.WithContext(ctx).Save(&model).Error; err != nil {
		return domain.User{}, err
	}
	return model.toDomain(), nil
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

// Record 持久化关键业务操作的审计记录。
func (s *GormStore) Record(ctx context.Context, entry audit.Entry) error {
	createdAt := entry.CreatedAt
	if createdAt.IsZero() {
		createdAt = time.Now()
	}
	return s.db.WithContext(ctx).Create(&auditModel{
		Action:    entry.Action,
		ActorID:   entry.ActorID,
		Resource:  entry.Resource,
		CreatedAt: createdAt,
	}).Error
}

// List 按写入顺序返回审计记录，供管理端分页查询。
func (s *GormStore) List(ctx context.Context) ([]audit.Entry, error) {
	var models []auditModel
	if err := s.db.WithContext(ctx).Order("id ASC").Find(&models).Error; err != nil {
		return nil, err
	}

	entries := make([]audit.Entry, 0, len(models))
	for _, model := range models {
		entries = append(entries, model.toAuditEntry())
	}
	return entries, nil
}

func (m userModel) toDomain() domain.User {
	return domain.User{
		ID:           m.ID,
		Name:         m.Name,
		Email:        m.Email,
		Role:         m.Role,
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

func (m auditModel) toAuditEntry() audit.Entry {
	return audit.Entry{
		Action:    m.Action,
		ActorID:   m.ActorID,
		Resource:  m.Resource,
		CreatedAt: m.CreatedAt,
	}
}

func mapGormError(err error) error {
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrNotFound
	}
	return err
}
