package repository

import (
	"context"
	"errors"
	"sync"
	"time"

	"GoStudy/internal/domain"
)

var ErrNotFound = errors.New("not found")

type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUser(ctx context.Context, id int64) (domain.User, error)
	FindUserByEmail(ctx context.Context, email string) (domain.User, error)
	ListUsers(ctx context.Context) ([]domain.User, error)
}

type TaskRepository interface {
	CreateTask(ctx context.Context, task domain.Task) (domain.Task, error)
	GetTask(ctx context.Context, id int64) (domain.Task, error)
	ListTasks(ctx context.Context) ([]domain.Task, error)
	UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error)
}

// MemoryStore 用内存模拟持久化，便于企业应用骨架先跑起来。
type MemoryStore struct {
	mu         sync.RWMutex
	nextUserID int64
	nextTaskID int64
	users      map[int64]domain.User
	tasks      map[int64]domain.Task
}

// NewMemoryStore 创建线程安全的内存仓储，适合本地示例和单元测试。
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		nextUserID: 1,
		nextTaskID: 1,
		users:      make(map[int64]domain.User),
		tasks:      make(map[int64]domain.Task),
	}
}

// CreateUser 创建用户并分配自增 ID。
func (s *MemoryStore) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	user.ID = s.nextUserID
	user.CreatedAt = time.Now()
	s.nextUserID++
	s.users[user.ID] = user
	return user, nil
}

// GetUser 按 ID 查询用户。
func (s *MemoryStore) GetUser(ctx context.Context, id int64) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	user, ok := s.users[id]
	if !ok {
		return domain.User{}, ErrNotFound
	}
	return user, nil
}

// FindUserByEmail 支持注册去重和登录查询。
func (s *MemoryStore) FindUserByEmail(ctx context.Context, email string) (domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, user := range s.users {
		if user.Email == email {
			return user, nil
		}
	}
	return domain.User{}, ErrNotFound
}

// ListUsers 返回当前内存中的用户快照。
func (s *MemoryStore) ListUsers(ctx context.Context) ([]domain.User, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	users := make([]domain.User, 0, len(s.users))
	for _, user := range s.users {
		users = append(users, user)
	}
	return users, nil
}

// CreateTask 创建任务并分配自增 ID。
func (s *MemoryStore) CreateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	now := time.Now()
	task.ID = s.nextTaskID
	task.CreatedAt = now
	task.UpdatedAt = now
	s.nextTaskID++
	s.tasks[task.ID] = task
	return task, nil
}

// GetTask 按 ID 查询任务。
func (s *MemoryStore) GetTask(ctx context.Context, id int64) (domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	task, ok := s.tasks[id]
	if !ok {
		return domain.Task{}, ErrNotFound
	}
	return task, nil
}

// ListTasks 返回当前内存中的任务快照。
func (s *MemoryStore) ListTasks(ctx context.Context) ([]domain.Task, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	tasks := make([]domain.Task, 0, len(s.tasks))
	for _, task := range s.tasks {
		tasks = append(tasks, task)
	}
	return tasks, nil
}

// UpdateTask 更新已存在的任务。
func (s *MemoryStore) UpdateTask(ctx context.Context, task domain.Task) (domain.Task, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, ok := s.tasks[task.ID]; !ok {
		return domain.Task{}, ErrNotFound
	}
	task.UpdatedAt = time.Now()
	s.tasks[task.ID] = task
	return task, nil
}
