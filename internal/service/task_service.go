package service

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"GoStudy/internal/domain"
	"GoStudy/internal/event"
	"GoStudy/internal/repository"
)

type CreateTaskInput struct {
	Title string `json:"title" binding:"required"`
}

type TaskService struct {
	tasks  repository.TaskRepository
	users  repository.UserRepository
	events event.Publisher
}

// NewTaskService 依赖任务仓储和用户仓储，用于校验任务归属。
func NewTaskService(tasks repository.TaskRepository, users repository.UserRepository, publishers ...event.Publisher) *TaskService {
	publisher := event.Publisher(event.NewMemory())
	if len(publishers) > 0 && publishers[0] != nil {
		publisher = publishers[0]
	}
	return &TaskService{tasks: tasks, users: users, events: publisher}
}

// Create 创建任务前会确认负责人存在，避免产生孤立任务。
func (s *TaskService) Create(ctx context.Context, ownerID int64, input CreateTaskInput) (domain.Task, error) {
	title := strings.TrimSpace(input.Title)
	if title == "" || ownerID <= 0 {
		return domain.Task{}, fmt.Errorf("%w: task title is required", ErrInvalidInput)
	}

	if _, err := s.users.GetUser(ctx, ownerID); err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Task{}, fmt.Errorf("%w: owner not found", ErrNotFound)
		}
		return domain.Task{}, err
	}

	task, err := s.tasks.CreateTask(ctx, domain.Task{
		Title:   title,
		OwnerID: ownerID,
		Status:  domain.TaskStatusPending,
	})
	if err != nil {
		return domain.Task{}, err
	}
	if err := s.events.Publish(ctx, "task.created", task); err != nil {
		return domain.Task{}, err
	}
	return task, nil
}

// TaskPage 是任务列表及其分页信息。
type TaskPage struct {
	PageMeta
	Items []domain.Task `json:"items"`
}

// List 返回当前用户拥有的分页任务列表，避免跨用户读取任务。
func (s *TaskService) List(ctx context.Context, ownerID int64, input PageInput) (TaskPage, error) {
	if ownerID <= 0 {
		return TaskPage{}, fmt.Errorf("%w: owner is required", ErrInvalidInput)
	}
	tasks, err := s.tasks.ListTasks(ctx)
	if err != nil {
		return TaskPage{}, err
	}
	owned := make([]domain.Task, 0, len(tasks))
	for _, task := range tasks {
		if task.OwnerID == ownerID {
			owned = append(owned, task)
		}
	}
	meta, start, end, err := normalizePage(input, len(owned))
	if err != nil {
		return TaskPage{}, err
	}
	return TaskPage{PageMeta: meta, Items: owned[start:end]}, nil
}

// Complete 仅允许任务所有者将任务状态流转为完成。
func (s *TaskService) Complete(ctx context.Context, ownerID int64, id int64) (domain.Task, error) {
	task, err := s.tasks.GetTask(ctx, id)
	if err != nil {
		if errors.Is(err, repository.ErrNotFound) {
			return domain.Task{}, fmt.Errorf("%w: task not found", ErrNotFound)
		}
		return domain.Task{}, err
	}
	if task.OwnerID != ownerID {
		return domain.Task{}, fmt.Errorf("%w: task does not belong to current user", ErrForbidden)
	}

	task.Status = domain.TaskStatusDone
	task, err = s.tasks.UpdateTask(ctx, task)
	if err != nil {
		return domain.Task{}, err
	}
	if err := s.events.Publish(ctx, "task.completed", task); err != nil {
		return domain.Task{}, err
	}
	return task, nil
}
