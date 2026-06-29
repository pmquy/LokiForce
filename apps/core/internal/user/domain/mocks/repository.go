package mocks

import "lokiforce.com/apps/core/internal/user/domain"

type mockUserRepository struct {
	users map[string]*domain.User
}

func (m *mockUserRepository) CreateUser(user *domain.User) error {
	if m.users == nil {
		m.users = make(map[string]*domain.User)
	}
	m.users[user.ID] = user
	return nil
}

func (m *mockUserRepository) GetUserByEmail(email string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Email == domain.Email(email) {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) GetUserByID(id string) (*domain.User, error) {
	if user, exists := m.users[id]; exists {
		return user, nil
	}
	return nil, domain.ErrUserNotFound
}

func (m *mockUserRepository) UpdateUser(user *domain.User) error {
	if _, exists := m.users[user.ID]; exists {
		m.users[user.ID] = user
		return nil
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) DeleteUser(id string) error {
	if _, exists := m.users[id]; exists {
		delete(m.users, id)
		return nil
	}
	return domain.ErrUserNotFound
}

func (m *mockUserRepository) ListUsers() ([]*domain.User, error) {
	var userList []*domain.User
	for _, user := range m.users {
		userList = append(userList, user)
	}
	return userList, nil
}

func (m *mockUserRepository) GetUserByUsername(username string) (*domain.User, error) {
	for _, user := range m.users {
		if user.Username == username {
			return user, nil
		}
	}
	return nil, domain.ErrUserNotFound
}

func NewMockUserRepository() *mockUserRepository {
	return &mockUserRepository{
		users: make(map[string]*domain.User),
	}
}
