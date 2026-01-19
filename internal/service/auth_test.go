package service

import (
	"context"
	"testing"
	"time"

	"golang.org/x/crypto/bcrypt"

	"bench-hub/internal/model"
)

func TestAuthServiceLoginRefresh(t *testing.T) {
	repo := newFakeUserRepo()
	hash, err := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	if err != nil {
		t.Fatalf("hash: %v", err)
	}
	_ = repo.Create(context.Background(), &model.User{
		Username:     "alice",
		PasswordHash: string(hash),
	})

	svc := NewAuthService(repo, "test-secret", 10*time.Minute, 7*24*time.Hour, "test-issuer")

	access, refresh, user, err := svc.Login(context.Background(), "alice", "secret")
	if err != nil {
		t.Fatalf("login: %v", err)
	}
	if access == "" || refresh == "" {
		t.Fatalf("expected tokens")
	}
	if user.Username != "alice" {
		t.Fatalf("expected user alice")
	}

	userID, err := svc.ValidateAccess(access)
	if err != nil {
		t.Fatalf("validate: %v", err)
	}
	if userID == "" {
		t.Fatalf("expected user id")
	}

	newAccess, newRefresh, err := svc.Refresh(context.Background(), refresh)
	if err != nil {
		t.Fatalf("refresh: %v", err)
	}
	if newAccess == "" || newRefresh == "" {
		t.Fatalf("expected refreshed tokens")
	}
}

func TestAuthServiceInvalidCredentials(t *testing.T) {
	repo := newFakeUserRepo()
	hash, _ := bcrypt.GenerateFromPassword([]byte("secret"), bcrypt.DefaultCost)
	_ = repo.Create(context.Background(), &model.User{
		Username:     "alice",
		PasswordHash: string(hash),
	})

	svc := NewAuthService(repo, "test-secret", 10*time.Minute, 7*24*time.Hour, "test-issuer")
	_, _, _, err := svc.Login(context.Background(), "alice", "wrong")
	if err != ErrInvalidCredentials {
		t.Fatalf("expected invalid credentials error")
	}
}
