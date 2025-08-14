package postgres

import (
	"context"
	"potapov.me/taskflow/internal/domain/project"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

func (r *Repository) CreateProject(ctx context.Context, p *project.Project) error {
	return r.DB.WithContext(ctx).Create(p).Error
}

func (r *Repository) GetProject(ctx context.Context, id uuid.UUID) (*project.Project, error) {
	var p project.Project
	err := r.DB.WithContext(ctx).First(&p, "id = ?", id).Error
	return &p, err
}

func (r *Repository) ListProjects(ctx context.Context, ownerID uuid.UUID) ([]*project.Project, error) {
	var projects []*project.Project
	err := r.DB.WithContext(ctx).Where("owner_id = ?", ownerID).Find(&projects).Error
	return projects, err
}

func (r *Repository) UpdateProject(ctx context.Context, id uuid.UUID, updateFn func(*project.Project) (*project.Project, error)) error {
	return r.DB.Transaction(func(tx *gorm.DB) error {
		var p project.Project
		if err := tx.WithContext(ctx).First(&p, "id = ?", id).Error; err != nil {
			return err
		}

		updated, err := updateFn(&p)
		if err != nil {
			return err
		}

		return tx.WithContext(ctx).Save(updated).Error
	})
}

func (r *Repository) DeleteProject(ctx context.Context, id uuid.UUID) error {
	return r.DB.WithContext(ctx).Delete(&project.Project{}, "id = ?", id).Error
}
