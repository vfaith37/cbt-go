package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	Email     string    `json:"email" gorm:"unique"`
	Password  string    `json:"-"`
	Role      string    `json:"role"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Exam struct {
	ID        uuid.UUID  `json:"id" gorm:"primaryKey;type:uuid"`
	Title     string     `json:"title"`
	Duration  int        `json:"duration"` // in minutes
	StartTime time.Time  `json:"start_time"`
	EndTime   time.Time  `json:"end_time"`
	Questions []Question `json:"questions"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
}

type Question struct {
	ID            uuid.UUID `json:"id" gorm:"primaryKey;type:uuid"`
	ExamID        uuid.UUID `json:"exam_id"`
	Content       string    `json:"content"`
	Options       []string  `json:"options" gorm:"type:jsonb"`
	CorrectAnswer string    `json:"correct_answer"`
	Marks         int       `json:"marks"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}
