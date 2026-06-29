package repository

import (
	"errors"

	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/project/domain"
)

type projectModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	OrgID       string `gorm:"not null;index"`
}

func (projectModel) TableName() string {
	return "projects"
}

type PostgresProjectRepository struct {
	db *gorm.DB
}

func NewPostgresProjectRepository(db *gorm.DB) *PostgresProjectRepository {
	return &PostgresProjectRepository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&projectModel{})
}

func (r *PostgresProjectRepository) Create(project *domain.Project) error {
	m := toModel(project)
	return r.db.Create(m).Error
}

func (r *PostgresProjectRepository) GetByID(id string) (*domain.Project, error) {
	var m projectModel
	if err := r.db.First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("project not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresProjectRepository) ListByOrg(orgID string) ([]*domain.Project, error) {
	var models []projectModel
	if err := r.db.Find(&models, "org_id = ?", orgID).Error; err != nil {
		return nil, err
	}

	projects := make([]*domain.Project, len(models))
	for i, m := range models {
		projects[i] = toDomain(&m)
	}
	return projects, nil
}

func (r *PostgresProjectRepository) Delete(id string) error {
	res := r.db.Delete(&projectModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("project not found")
	}
	return nil
}

func toModel(p *domain.Project) *projectModel {
	return &projectModel{
		ID:          p.ID,
		Name:        p.Name,
		Description: p.Description,
		OrgID:       p.OrgID,
	}
}

func toDomain(m *projectModel) *domain.Project {
	return &domain.Project{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		OrgID:       m.OrgID,
	}
}
