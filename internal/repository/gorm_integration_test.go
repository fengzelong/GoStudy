package repository

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"GoStudy/internal/domain"
)

func TestGormStoreIntegration(t *testing.T) {
	if os.Getenv("APP_INTEGRATION_MYSQL") != "1" {
		t.Skip("未设置 APP_INTEGRATION_MYSQL=1，跳过企业骨架 MySQL 集成测试")
	}
	dsn := os.Getenv("MYSQL_DSN")
	if dsn == "" {
		t.Skip("未设置 MYSQL_DSN，跳过企业骨架 MySQL 集成测试")
	}

	store, err := NewGormStore(dsn)
	if err != nil {
		t.Fatalf("创建 GORM 仓储失败: %v", err)
	}
	defer store.Close()

	ctx := context.Background()
	suffix := time.Now().UnixNano()
	user, err := store.CreateUser(ctx, domain.User{
		Name:         "integration-user",
		Email:        fmt.Sprintf("integration-%d@example.com", suffix),
		Role:         domain.UserRoleUser,
		PasswordHash: "integration-password-hash",
	})
	if err != nil {
		t.Fatalf("创建用户失败: %v", err)
	}
	defer func() {
		_ = store.db.WithContext(ctx).Delete(&taskModel{}, "owner_id = ?", user.ID).Error
		_ = store.db.WithContext(ctx).Delete(&userModel{}, user.ID).Error
	}()

	user.Role = domain.UserRoleAdmin
	user, err = store.UpdateUser(ctx, user)
	if err != nil {
		t.Fatalf("更新用户角色失败: %v", err)
	}
	if user.Role != domain.UserRoleAdmin {
		t.Fatalf("期望管理员角色，实际为 %s", user.Role)
	}

	byEmail, err := store.FindUserByEmail(ctx, user.Email)
	if err != nil || byEmail.ID != user.ID {
		t.Fatalf("按邮箱查询用户失败: user=%+v err=%v", byEmail, err)
	}

	task, err := store.CreateTask(ctx, domain.Task{
		Title:   "integration-task",
		OwnerID: user.ID,
		Status:  domain.TaskStatusPending,
	})
	if err != nil {
		t.Fatalf("创建任务失败: %v", err)
	}
	task.Status = domain.TaskStatusDone
	task, err = store.UpdateTask(ctx, task)
	if err != nil {
		t.Fatalf("更新任务失败: %v", err)
	}
	if task.Status != domain.TaskStatusDone {
		t.Fatalf("期望完成状态，实际为 %s", task.Status)
	}

	storedTask, err := store.GetTask(ctx, task.ID)
	if err != nil || storedTask.OwnerID != user.ID || storedTask.Status != domain.TaskStatusDone {
		t.Fatalf("查询任务失败: task=%+v err=%v", storedTask, err)
	}
}
