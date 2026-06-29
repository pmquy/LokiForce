package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/user/domain"
)

type userModel struct {
	ID       string `gorm:"primaryKey"`
	Username string `gorm:"uniqueIndex;not null"`
	Email    string `gorm:"uniqueIndex;not null"`
	Password string `gorm:"not null"`
	Role     string `gorm:"not null"`
}

func (userModel) TableName() string {
	return "users"
}

type PostgresUserRepository struct {
	db *gorm.DB
}

func NewPostgresUserRepository(db *gorm.DB) *PostgresUserRepository {
	return &PostgresUserRepository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&userModel{})
}

func (r *PostgresUserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	if err := r.db.WithContext(ctx).Create(m).Error; err != nil {
		return err
	}
	return nil
}

func (r *PostgresUserRepository) GetUserByID(ctx context.Context, id string) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresUserRepository) GetUserByUsername(ctx context.Context, username string) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).First(&m, "username = ?", username).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresUserRepository) GetUserByEmail(ctx context.Context, email string) (*domain.User, error) {
	var m userModel
	if err := r.db.WithContext(ctx).First(&m, "email = ?", email).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, domain.ErrUserNotFound
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresUserRepository) UpdateUser(ctx context.Context, user *domain.User) error {
	m := toModel(user)
	res := r.db.WithContext(ctx).Save(m)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func (r *PostgresUserRepository) DeleteUser(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&userModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return domain.ErrUserNotFound
	}
	return nil
}

func toModel(u *domain.User) *userModel {
	return &userModel{
		ID:       u.ID,
		Username: u.Username,
		Email:    string(u.Email),
		Password: string(u.Password),
		Role:     string(u.Role),
	}
}

func toDomain(m *userModel) *domain.User {
	return &domain.User{
		ID:       m.ID,
		Username: m.Username,
		Email:    domain.Email(m.Email),
		Password: m.Password,
		Role:     domain.Role(m.Role),
	}
}
