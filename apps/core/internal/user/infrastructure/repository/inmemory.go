package repository

import (
	"context"
	"sync"

	"lokiforce.com/apps/core/internal/user/domain"
)

type InMemoryUserRepository struct {
	mu    sync.RWMutex
	users map[string]*domain.User
}

func NewInMemoryUserRepository() *InMemoryUserRepository {
	return &InMemoryUserRepository{
		users: make(map[string]*domain.User),
	}
}

func (r *InMemoryUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	for _, u := range r.users {
		if u.Email == user.Email {
			return domain.ErrInvalidEmail
		}
		if u.Username == user.Username {
			return domain.ErrInvalidCredentials
		}
	}

	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	user, exists := r.users[id]
	if !exists {
		return nil, domain.ErrUserNotFound
	}
	return user, nil
}

func (r *InMemoryUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *InMemoryUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	for _, user := range r.users {
		if string(user.Email) == email {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (r *InMemoryUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[user.ID]; !exists {
		return domain.ErrUserNotFound
	}
	r.users[user.ID] = user
	return nil
}

func (r *InMemoryUserRepository) DeleteUser(ctx context.Context, id string) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	if _, exists := r.users[id]; !exists {
		return domain.ErrUserNotFound
	}
	delete(r.users, id)
	return nil
}
