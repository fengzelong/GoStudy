package service

import (
	"context"
	"errors"
	"testing"

	"GoStudy/internal/audit"
	"GoStudy/internal/domain"
	"GoStudy/internal/event"
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

	task, err := taskSvc.Create(ctx, user.ID, CreateTaskInput{Title: "发布企业骨架"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if task.Status != domain.TaskStatusPending {
		t.Fatalf("expected pending task, got %s", task.Status)
	}

	task, err = taskSvc.Complete(ctx, user.ID, task.ID)
	if err != nil {
		t.Fatalf("complete task: %v", err)
	}
	if task.Status != domain.TaskStatusDone {
		t.Fatalf("expected done task, got %s", task.Status)
	}
}

func TestTaskServiceOwnershipAndPagination(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	userSvc := NewUserService(store)
	taskSvc := NewTaskService(store, store)

	alice, err := userSvc.Register(ctx, RegisterUserInput{Name: "Alice", Email: "alice@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register Alice: %v", err)
	}
	bob, err := userSvc.Register(ctx, RegisterUserInput{Name: "Bob", Email: "bob@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register Bob: %v", err)
	}
	for _, title := range []string{"A1", "A2", "A3"} {
		if _, err := taskSvc.Create(ctx, alice.ID, CreateTaskInput{Title: title}); err != nil {
			t.Fatalf("create Alice task: %v", err)
		}
	}
	bobTask, err := taskSvc.Create(ctx, bob.ID, CreateTaskInput{Title: "B1"})
	if err != nil {
		t.Fatalf("create Bob task: %v", err)
	}

	page, err := taskSvc.List(ctx, alice.ID, PageInput{Page: 2, PageSize: 2})
	if err != nil {
		t.Fatalf("list Alice tasks: %v", err)
	}
	if page.Total != 3 || len(page.Items) != 1 || page.Items[0].Title != "A3" {
		t.Fatalf("unexpected Alice page: %+v", page)
	}
	if _, err := taskSvc.Complete(ctx, alice.ID, bobTask.ID); !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected forbidden error, got %v", err)
	}
	if _, err := taskSvc.List(ctx, alice.ID, PageInput{Page: 0, PageSize: 101}); !errors.Is(err, ErrInvalidInput) {
		t.Fatalf("expected invalid pagination error, got %v", err)
	}
}

func TestTaskServicePublishesEvents(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	publisher := event.NewMemory()
	userSvc := NewUserService(store)
	taskSvc := NewTaskService(store, store, publisher)
	user, err := userSvc.Register(ctx, RegisterUserInput{Name: "Alice", Email: "alice@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}
	task, err := taskSvc.Create(ctx, user.ID, CreateTaskInput{Title: "publish event"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if _, err := taskSvc.Complete(ctx, user.ID, task.ID); err != nil {
		t.Fatalf("complete task: %v", err)
	}
	messages := publisher.Messages()
	if len(messages) != 2 || messages[0].Name != "task.created" || messages[1].Name != "task.completed" {
		t.Fatalf("unexpected task events: %+v", messages)
	}
}

func TestServicesRecordAuditEntries(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	logger := audit.NewMemoryLogger()
	userSvc := NewUserService(store)
	userSvc.SetAuditLogger(logger)
	taskSvc := NewTaskService(store, store)
	taskSvc.SetAuditLogger(logger)
	user, err := userSvc.Register(ctx, RegisterUserInput{Name: "Alice", Email: "alice@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register user: %v", err)
	}
	if _, err := userSvc.Login(ctx, LoginInput{Email: user.Email, Password: "secret1"}); err != nil {
		t.Fatalf("login user: %v", err)
	}
	task, err := taskSvc.Create(ctx, user.ID, CreateTaskInput{Title: "audit task"})
	if err != nil {
		t.Fatalf("create task: %v", err)
	}
	if _, err := taskSvc.Complete(ctx, user.ID, task.ID); err != nil {
		t.Fatalf("complete task: %v", err)
	}
	entries := logger.Entries()
	if len(entries) != 3 || entries[0].Action != "auth.login" || entries[1].Action != "task.created" || entries[2].Action != "task.completed" {
		t.Fatalf("unexpected audit entries: %+v", entries)
	}
}

func TestUserServiceListPagination(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	svc := NewUserService(store)
	for _, email := range []string{"a@example.com", "b@example.com", "c@example.com"} {
		if _, err := svc.Register(ctx, RegisterUserInput{Name: email, Email: email, Password: "secret1"}); err != nil {
			t.Fatalf("register user: %v", err)
		}
	}
	page, err := svc.List(ctx, PageInput{Page: 2, PageSize: 2})
	if err != nil {
		t.Fatalf("list users: %v", err)
	}
	if page.Page != 2 || page.PageSize != 2 || page.Total != 3 || len(page.Items) != 1 {
		t.Fatalf("unexpected user page: %+v", page)
	}
}

func TestUserServiceAssignsConfiguredAdminRole(t *testing.T) {
	ctx := context.Background()
	store := repository.NewMemoryStore()
	svc := NewUserService(store)
	svc.SetAdminEmail("admin@example.com")
	admin, err := svc.Register(ctx, RegisterUserInput{Name: "Admin", Email: "admin@example.com", Password: "secret1"})
	if err != nil {
		t.Fatalf("register admin: %v", err)
	}
	if admin.Role != domain.UserRoleAdmin {
		t.Fatalf("expected admin role, got %s", admin.Role)
	}
}
