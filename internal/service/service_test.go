package service

import (
	"context"
	"errors"
	"testing"

	"GoStudy/internal/domain"
	"GoStudy/internal/repository"
)

func TestUserServiceRegister(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	svc := NewUserService(store)

	user, err := svc.Register(ctx, RegisterUserInput{Name: "Alice", Email: "Alice@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}
	if user.ID == 0 {
		t.Fatal("expected user id")
	}
	if user.Email != "alice@example.com" {
		t.Fatalf("expected normalized email, got %s", user.Email)
	}

	_, err = svc.Register(ctx, RegisterUserInput{Name: "Alice2", Email: "alice@example.com", Password: "secret1"})
	if !errors.Is(err, ErrConflict) {
		t.Fatalf("expected conflict, got %v", err)
	}

	_, err = svc.Login(ctx, LoginInput{Email: "alice@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("login user: %v", err)
	}
}

func TestTaskServiceComplete(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	userSvc := NewUserService(store)
	taskSvc := NewTaskService(store, store)

	user, err := userSvc.Register(ctx, RegisterUserInput{Name: "Bob", Email: "bob@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}

	task, err := taskSvc.Create(ctx, CreateTaskInput{Title: "发布企业骨架", OwnerID: user.ID})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if task.Status != domain.TaskStatusPending {
		t.Fatalf("expected pending task, got %s", task.Status)
	}

	task, err = taskSvc.Complete(ctx, task.ID)
	if err != nil {
		t.Fatalf("complete task: %v", err)
	}
	if task.Status != domain.TaskStatusDone {
		t.Fatalf("expected done task, got %s", task.Status)
	}
}
