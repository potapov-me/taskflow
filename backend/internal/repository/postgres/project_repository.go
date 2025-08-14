package postgres

import (
	"context"
	"potapov.me/taskflow/internal/domain/project"
)

func (r *Repository) CreateProject(ctx context.Context, p *project.Project) error {
	query := `INSERT INTO projects (id, title, description, owner_id) 
              VALUES ($1, $2, $3, $4)`
	_, err := r.Pool.Exec(ctx, query, p.ID, p.Title, p.Description, p.OwnerID)
	return err
}

// Добавьте методы GetByID, Update, Delete, List...
