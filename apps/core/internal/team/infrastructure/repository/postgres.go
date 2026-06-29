package repository

import (
	"context"
	"errors"

	"gorm.io/gorm"
	"lokiforce.com/apps/core/internal/team/domain"
)

type teamModel struct {
	ID          string `gorm:"primaryKey"`
	Name        string `gorm:"not null"`
	Description string
	OrgID       string `gorm:"not null;index"`
}

func (teamModel) TableName() string {
	return "teams"
}

type PostgresTeamRepository struct {
	db *gorm.DB
}

func NewPostgresTeamRepository(db *gorm.DB) *PostgresTeamRepository {
	return &PostgresTeamRepository{db: db}
}

func Migrate(db *gorm.DB) error {
	return db.AutoMigrate(&teamModel{})
}

func (r *PostgresTeamRepository) Create(ctx context.Context, team *domain.Team) error {
	m := toModel(team)
	return r.db.WithContext(ctx).Create(m).Error
}

func (r *PostgresTeamRepository) GetByID(ctx context.Context, id string) (*domain.Team, error) {
	var m teamModel
	if err := r.db.WithContext(ctx).First(&m, "id = ?", id).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("team not found")
		}
		return nil, err
	}
	return toDomain(&m), nil
}

func (r *PostgresTeamRepository) ListByOrg(ctx context.Context, orgID string) ([]*domain.Team, error) {
	var models []teamModel
	if err := r.db.WithContext(ctx).Find(&models, "org_id = ?", orgID).Error; err != nil {
		return nil, err
	}

	teams := make([]*domain.Team, len(models))
	for i, m := range models {
		teams[i] = toDomain(&m)
	}
	return teams, nil
}

func (r *PostgresTeamRepository) Delete(ctx context.Context, id string) error {
	res := r.db.WithContext(ctx).Delete(&teamModel{}, "id = ?", id)
	if res.Error != nil {
		return res.Error
	}
	if res.RowsAffected == 0 {
		return errors.New("team not found")
	}
	return nil
}

func toModel(t *domain.Team) *teamModel {
	return &teamModel{
		ID:          t.ID,
		Name:        t.Name,
		Description: t.Description,
		OrgID:       t.OrgID,
	}
}

func toDomain(m *teamModel) *domain.Team {
	return &domain.Team{
		ID:          m.ID,
		Name:        m.Name,
		Description: m.Description,
		OrgID:       m.OrgID,
	}
}
