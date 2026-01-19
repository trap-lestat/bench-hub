package service

import (
	"context"
	"testing"
)

func TestUserServiceCRUD(t *testing.T) {
	repo := newFakeUserRepo()
	svc := NewUserService(repo)

	user, err := svc.Create(context.Background(), "alice", "password")
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if user.ID == "" {
		t.Fatalf("expected user id")
	}

	got, err := svc.Get(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	if got.Username != "alice" {
		t.Fatalf("expected username alice, got %s", got.Username)
	}

	updated, err := svc.Update(context.Background(), user.ID, "bob", "newpass")
	if err != nil {
		t.Fatalf("update: %v", err)
	}
	if updated.Username != "bob" {
		t.Fatalf("expected username bob, got %s", updated.Username)
	}

	if err := svc.Delete(context.Background(), user.ID); err != nil {
		t.Fatalf("delete: %v", err)
	}

	if _, err := svc.Get(context.Background(), user.ID); err != ErrNotFound {
		t.Fatalf("expected not found after delete")
	}
}
