package repository

import (
	"context"
	"time"

	"github.com/your-org/testing-platform/core/internal/models"
)

// Repository defines the interface for data storage operations
type Repository interface {
	// Task operations
	CreateTask(ctx context.Context, task *models.Task) error
	GetTask(ctx context.Context, taskID string) (*models.Task, error)
	UpdateTask(ctx context.Context, task *models.Task) error
	DeleteTask(ctx context.Context, taskID string) error
	ListTasks(ctx context.Context, filter TaskFilter) ([]models.Task, error)
	SearchTasks(ctx context.Context, query string) ([]models.Task, error)

	// Session operations
	CreateSession(ctx context.Context, session *models.TestSession) error
	GetSession(ctx context.Context, sessionID string) (*models.TestSession, error)
	UpdateSession(ctx context.Context, session *models.TestSession) error
	ListActiveSessions(ctx context.Context) ([]models.TestSession, error)
	CleanupExpiredSessions(ctx context.Context, before time.Time) error

	// Answer operations
	SaveAnswer(ctx context.Context, answer *models.Answer) error
	GetAnswers(ctx context.Context, sessionID string) ([]models.Answer, error)

	// Result operations
	SaveResult(ctx context.Context, result *models.TestResult) error
	GetResult(ctx context.Context, sessionID string) (*models.TestResult, error)
	GetResults(ctx context.Context, filter ResultFilter) ([]models.TestResult, error)

	// Analytics operations
	SaveAnalyticsReport(ctx context.Context, report *models.AnalyticsReport) error
	GetAnalyticsReport(ctx context.Context, reportID string) (*models.AnalyticsReport, error)
}

// TaskFilter defines filtering options for tasks
type TaskFilter struct {
	AuthorEmail string
	Difficulty  string
	Tags        []string
	IsPublic    *bool
	CreatedFrom *time.Time
	CreatedTo   *time.Time
}

// ResultFilter defines filtering options for test results
type ResultFilter struct {
	UserEmail  string
	TaskID     string
	StartDate  *time.Time
	EndDate    *time.Time
	MinScore   *float64
	MaxScore   *float64
	Status     string
}
