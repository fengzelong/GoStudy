package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"GoStudy/internal/domain"
	"GoStudy/internal/repository"
)

type CreateTaskInput struct {
	Title   string `json:"title" binding:"required"`
	OwnerID int64  `json:"owner_id" binding:"required"`
}

type TaskService struct {
	tasks repository.TaskRepository
	users repository.UserRepository
}

// NewTaskService 依赖任务仓储和用户仓储，用于校验任务归属。
func NewTaskService(tasks repository.TaskRepository, users repository.UserRepository) *TaskService {
	return &TaskService{tasks: tasks, users: users}
}

// Create 创建任务前会确认负责人存在，避免产生孤立任务。
func (s *TaskService) Create(ctx context.Context, input CreateTaskInput) (domain.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" || input.OwnerID <= 0 {
		return domain.Task{}, fmt.Errorf("%w: task title and owner_id are required", ErrInvalidInput)
	}

	if _, err := s.users.GetUser(ctx, input.OwnerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Task{}, fmt.Errorf("%w: owner not found", ErrNotFound)
		}
		return domain.Task{}, err
	}

	return s.tasks.CreateTask(ctx, domain.Task{
		Title:   title,
		OwnerID: input.OwnerID,
		Status:  domain.TaskStatusPending,
	})
}

// List 返回任务列表，排序规则由仓储实现保证。
func (s *TaskService) List(ctx context.Context) ([]domain.Task, error) {
	return s.tasks.ListTasks(ctx)
}

// Complete 将任务状态流转为完成。
func (s *TaskService) Complete(ctx context.Context, id int64) (domain.Task, error) {
	task, err := s.tasks.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Task{}, fmt.Errorf("%w: task not found", ErrNotFound)
		}
		return domain.Task{}, err
	}

	task.Status = domain.TaskStatusDone
	return s.tasks.UpdateTask(ctx, task)
}
