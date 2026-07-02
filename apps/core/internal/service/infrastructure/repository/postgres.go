package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/service/domain"
)

type serviceModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	ProjectID   string `gorm:"not null;index"`
	TemplateID  string `gorm:"not null"`
	Repository  string
}

func (serviceModel) TableName() string {
	return "services"
}

type accessPolicyModel struct {
	ID         string `gorm:"primaryKey"`
	ClientID   string `gorm:"not null"`
	TargetID   string `gorm:"not null"`
	TargetPort string `gorm:"not null"`
	ProjectID  string `gorm:"not null;index"`
}

func (accessPolicyModel) TableName() string {
	return "access_policies"
}

type PostgresServiceRepository struct {
	db *gorm.DB
}

func NewPostgresServiceRepository(db *gorm.DB) *PostgresServiceRepository {
	return &PostgresServiceRepository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&serviceModel{}, &accessPolicyModel{})
}

func (r *PostgresServiceRepository) Create(ctx context.Context, service *domain.Service) error {
	m := toModel(service)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PostgresServiceRepository) GetByID(ctx context.Context, id string) (*domain.Service, error) {
	var m serviceModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("service not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresServiceRepository) ListByProject(ctx context.Context, projectID string) ([]*domain.Service, error) {
	var models []serviceModel
	if err := r.db.WithContext(ctx).Find(&models, "project_id = ?", projectID).Error; err != nil {
		return nil, err
	}

	services := make([]*domain.Service, len(models))
	for i, m := range models {
		services[i] = toDomain(&m)
	}
	return services, nil
}

func (r *PostgresServiceRepository) Update(ctx context.Context, service *domain.Service) error {
	m := toModel(service)
	res := r.db.WithContext(ctx).Save(m)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("service not found")
	}
	return nil
}

func (r *PostgresServiceRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&serviceModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("service not found")
	}
	return nil
}

func toModel(s *domain.Service) *serviceModel {
	return &serviceModel{
		ID:          s.ID,
		Name:        s.Name,
		Description: s.Description,
		ProjectID:   s.ProjectID,
		TemplateID:  s.TemplateID,
		Repository:  s.Repository,
	}
}

func toDomain(m *serviceModel) *domain.Service {
	return &domain.Service{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		ProjectID:   m.ProjectID,
		TemplateID:  m.TemplateID,
		Repository:  m.Repository,
	}
}

func (r *PostgresServiceRepository) CreateAccessPolicy(ctx context.Context, policy *domain.AccessPolicy) error {
	policyModel := &accessPolicyModel{
		ID:         policy.ID,
		ClientID:   policy.ClientID,
		TargetID:   policy.TargetID,
		TargetPort: policy.TargetPort,
		ProjectID:  policy.ProjectID,
	}
	return r.db.WithContext(ctx).Create(policyModel).Error
}

func (r *PostgresServiceRepository) DeleteAccessPolicy(ctx context.Context, policyID string) error {
	res := r.db.WithContext(ctx).Delete(&accessPolicyModel{}, "id = ?", policyID)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("access policy not found")
	}
	return nil
}

func (r *PostgresServiceRepository) ListAccessPoliciesByTarget(ctx context.Context, targetID string) ([]*domain.AccessPolicy, error) {
	var models []accessPolicyModel
	if err := r.db.WithContext(ctx).Find(&models, "target_id = ?", targetID).Error; err != nil {
		return nil, err
	}
	policies := make([]*domain.AccessPolicy, len(models))
	for i, m := range models {
		policies[i] = &domain.AccessPolicy{
			ID:         m.ID,
			ClientID:   m.ClientID,
			TargetID:   m.TargetID,
			TargetPort: m.TargetPort,
			ProjectID:  m.ProjectID,
		}
	}
	return policies, nil
}

func (r *PostgresServiceRepository) GetAccessPolicyByID(ctx context.Context, id string) (*domain.AccessPolicy, error) {
	var m accessPolicyModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("access policy not found")
		}
		return nil, err
	}
	return &domain.AccessPolicy{
		ID:         m.ID,
		ClientID:   m.ClientID,
		TargetID:   m.TargetID,
		TargetPort: m.TargetPort,
		ProjectID:  m.ProjectID,
	}, nil
}
