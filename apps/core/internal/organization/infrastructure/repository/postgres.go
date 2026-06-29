package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/organization/domain"
)

type orgModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	OwnerID     string `gorm:"not null;index"`
}

func (orgModel) TableName() string {
	return "organizations"
}

type PostgresOrgRepository struct {
	db *gorm.DB
}

func NewPostgresOrgRepository(db *gorm.DB) *PostgresOrgRepository {
	return &PostgresOrgRepository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&orgModel{})
}

func (r *PostgresOrgRepository) Create(ctx context.Context, org *domain.Organization) error {
	m := toModel(org)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PostgresOrgRepository) GetByID(ctx context.Context, id string) (*domain.Organization, error) {
	var m orgModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("organization not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresOrgRepository) ListByOwner(ctx context.Context, ownerID string) ([]*domain.Organization, error) {
	var models []orgModel
	if err := r.db.WithContext(ctx).Find(&models, "owner_id = ?", ownerID).Error; err != nil {
		return nil, err
	}

	orgs := make([]*domain.Organization, len(models))
	for i, m := range models {
		orgs[i] = toDomain(&m)
	}
	return orgs, nil
}

func (r *PostgresOrgRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&orgModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("organization not found")
	}
	return nil
}

func toModel(org *domain.Organization) *orgModel {
	return &orgModel{
		ID:          org.ID,
		Name:        org.Name,
		Description: org.Description,
		OwnerID:     org.OwnerID,
	}
}

func toDomain(m *orgModel) *domain.Organization {
	return &domain.Organization{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		OwnerID:     m.OwnerID,
	}
}
