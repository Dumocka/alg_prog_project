package model

import (
	"time"
)

type User struct {
	ID        int64     `json:"id"`
	Email     string    `json:"email"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Test struct {
	ID          int64     `json:"id"`
	Title       string    `json:"title"`
	Description string    `json:"description"`
	CreatorID   int64     `json:"creator_id"`
	IsPublic    bool      `json:"is_public"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Question struct {
	ID       int64    `json:"id"`
	TestID   int64    `json:"test_id"`
	Type     string   `json:"type"` // multiple_choice, single_choice, text
	Text     string   `json:"text"`
	Options  []string `json:"options,omitempty"`
	Required bool     `json:"required"`
}

type Answer struct {
	ID         int64     `json:"id"`
	QuestionID int64     `json:"question_id"`
	UserID     int64     `json:"user_id"`
	Value      string    `json:"value"`
	CreatedAt  time.Time `json:"created_at"`
}

type TestResult struct {
	ID        int64     `json:"id"`
	TestID    int64     `json:"test_id"`
	UserID    int64     `json:"user_id"`
	Score     float64   `json:"score,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

type TestPermission struct {
	ID     int64  `json:"id"`
	TestID int64  `json:"test_id"`
	UserID int64  `json:"user_id"`
	Role   string `json:"role"` // owner, editor, viewer
}
