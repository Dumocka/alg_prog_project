package models

import (
	"time"
)

// Task представляет задачу в системе
type Task struct {
	ID          string    `json:"id" bson:"_id"`
	Title       string    `json:"title" bson:"title"`
	Description string    `json:"description" bson:"description"`
	Difficulty  string    `json:"difficulty" bson:"difficulty"`
	Tags        []string  `json:"tags" bson:"tags"`
	AuthorEmail string    `json:"author_email" bson:"author_email"`
	CreatedAt   time.Time `json:"created_at" bson:"created_at"`
	IsPublic    bool      `json:"is_public" bson:"is_public"`
}

// TaskCreationRequest запрос на создание задачи
type TaskCreationRequest struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Difficulty  string   `json:"difficulty"`
	Tags        []string `json:"tags"`
	IsPublic    bool     `json:"is_public"`
}

// TestSession представляет сессию тестирования
type TestSession struct {
	ID          string     `json:"id" bson:"_id"`
	UserEmail   string     `json:"user_email" bson:"user_email"`
	TaskID      string     `json:"task_id" bson:"task_id"`
	StartedAt   time.Time  `json:"started_at" bson:"started_at"`
	CompletedAt *time.Time `json:"completed_at,omitempty" bson:"completed_at,omitempty"`
	Status      string     `json:"status" bson:"status"` // "in_progress", "completed", "timeout"
}

// Answer представляет ответ пользователя
type Answer struct {
	SessionID   string    `json:"session_id" bson:"session_id"`
	QuestionID  string    `json:"question_id" bson:"question_id"`
	Content     string    `json:"content" bson:"content"`
	SubmittedAt time.Time `json:"submitted_at" bson:"submitted_at"`
}

// TestResult представляет результат тестирования
type TestResult struct {
	SessionID      string    `json:"session_id" bson:"_id"`
	UserEmail     string    `json:"user_email" bson:"user_email"`
	TaskID        string    `json:"task_id" bson:"task_id"`
	Score         int       `json:"score" bson:"score"`
	CorrectAnswers []string  `json:"correct_answers" bson:"correct_answers"`
	WrongAnswers   []string  `json:"wrong_answers" bson:"wrong_answers"`
	CompletedAt    time.Time `json:"completed_at" bson:"completed_at"`
}

// AnalyticsRequest запрос на получение аналитики
type AnalyticsRequest struct {
	UserEmail  *string    `json:"user_email,omitempty"`
	TaskID     *string    `json:"task_id,omitempty"`
	StartDate  *time.Time `json:"start_date,omitempty"`
	EndDate    *time.Time `json:"end_date,omitempty"`
	Metrics    []string   `json:"metrics"`
}

// AnalyticsReport отчет по аналитике
type AnalyticsReport struct {
	ReportID    string                 `json:"report_id" bson:"_id"`
	Metrics     map[string]interface{} `json:"metrics" bson:"metrics"`
	GeneratedAt time.Time             `json:"generated_at" bson:"generated_at"`
}

// SystemStatus статус системы
type SystemStatus struct {
	IsHealthy       bool     `json:"is_healthy"`
	Version         string   `json:"version"`
	ActiveSessions  int      `json:"active_sessions"`
	CPUUsage        float64  `json:"cpu_usage"`
	MemoryUsage     float64  `json:"memory_usage"`
	ActiveServices  []string `json:"active_services"`
}

// SystemSettings настройки системы
type SystemSettings struct {
	MaxConcurrentSessions int      `json:"max_concurrent_sessions" bson:"max_concurrent_sessions"`
	SessionTimeoutMinutes int      `json:"session_timeout_minutes" bson:"session_timeout_minutes"`
	MaintenanceMode       bool     `json:"maintenance_mode" bson:"maintenance_mode"`
	EnabledFeatures       []string `json:"enabled_features" bson:"enabled_features"`
}
