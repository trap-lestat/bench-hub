package service

import (
	"context"
	"strconv"
	"sync"
	"time"

	"bench-hub/internal/model"
	"bench-hub/internal/repository"
)

type fakeUserRepo struct {
	mu       sync.Mutex
	users    map[string]*model.User
	byName   map[string]string
	sequence int
}

func newFakeUserRepo() *fakeUserRepo {
	return &fakeUserRepo{
		users:  make(map[string]*model.User),
		byName: make(map[string]string),
	}
}

func (r *fakeUserRepo) Create(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	r.sequence++
	if user.ID == "" {
		user.ID = "u-" + strconv.Itoa(r.sequence)
	}
	user.CreatedAt = time.Now()
	clone := *user
	r.users[user.ID] = &clone
	r.byName[user.Username] = user.ID
	return nil
}

func (r *fakeUserRepo) GetByID(ctx context.Context, id string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return nil, repository.ErrNotFound
	}
	clone := *user
	return &clone, nil
}

func (r *fakeUserRepo) GetByUsername(ctx context.Context, username string) (*model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	id, ok := r.byName[username]
	if !ok {
		return nil, repository.ErrNotFound
	}
	user := r.users[id]
	clone := *user
	return &clone, nil
}

func (r *fakeUserRepo) List(ctx context.Context, limit, offset int) ([]model.User, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	var out []model.User
	for _, user := range r.users {
		out = append(out, *user)
	}
	return out, nil
}

func (r *fakeUserRepo) Update(ctx context.Context, user *model.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, ok := r.users[user.ID]; !ok {
		return repository.ErrNotFound
	}
	clone := *user
	r.users[user.ID] = &clone
	r.byName[user.Username] = user.ID
	return nil
}

func (r *fakeUserRepo) Delete(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	user, ok := r.users[id]
	if !ok {
		return repository.ErrNotFound
	}
	delete(r.byName, user.Username)
	delete(r.users, id)
	return nil
}

func (r *fakeUserRepo) Count(ctx context.Context) (int, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	return len(r.users), nil
}
