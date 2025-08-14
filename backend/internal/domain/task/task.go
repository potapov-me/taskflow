package task

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Status string

const (
	Todo       Status = "todo"
	InProgress Status = "in_progress"
	Done       Status = "done"
)

type Task struct {
	ID          uuid.UUID  `json:"id" gorm:"type:uuid;primaryKey;default:uuid_generate_v4()"`
	ProjectID   uuid.UUID  `json:"projectId" gorm:"type:uuid;not null;index"`
	Title       string     `json:"title" gorm:"not null"`
	Description string     `json:"description"`
	Status      Status     `json:"status" gorm:"type:varchar(20);default:'todo'"`
	DueDate     *time.Time `json:"dueDate"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}

func (t *Task) BeforeCreate(tx *gorm.DB) (err error) {
	t.ID = uuid.New()
	return
}
