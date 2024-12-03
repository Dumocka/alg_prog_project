package service

import (
	"context"
	"time"

	"github.com/your-org/testing-platform/core/internal/errors"
	"github.com/your-org/testing-platform/core/internal/models"
	"github.com/your-org/testing-platform/core/internal/repository"
)

// TestRunner управляет процессом тестирования
type TestRunner struct {
	repo repository.Repository
}

// NewTestRunner создает новый экземпляр TestRunner
func NewTestRunner(repo repository.Repository) *TestRunner {
	return &TestRunner{repo: repo}
}

// CreateSession создает новую сессию тестирования
func (r *TestRunner) CreateSession(ctx context.Context, userEmail, taskID string) (*models.TestSession, error) {
	if userEmail == "" {
		return nil, errors.ErrInvalidUserEmail
	}
	if taskID == "" {
		return nil, errors.ErrInvalidTaskID
	}

	session := &models.TestSession{
		UserEmail: userEmail,
		TaskID:    taskID,
		StartedAt: time.Now(),
		Status:    "in_progress",
	}

	if err := r.repo.CreateSession(ctx, session); err != nil {
		return nil, err
	}

	return session, nil
}

// GetSession получает информацию о сессии
func (r *TestRunner) GetSession(ctx context.Context, sessionID string) (*models.TestSession, error) {
	session, err := r.repo.GetSession(ctx, sessionID)
	if err != nil {
		return nil, errors.ErrSessionNotFound
	}
	return session, nil
}

// SubmitAnswer сохраняет ответ пользователя
func (r *TestRunner) SubmitAnswer(ctx context.Context, answer models.Answer) error {
	session, err := r.GetSession(ctx, answer.SessionID)
	if err != nil {
		return err
	}

	if session.Status != "in_progress" {
		return errors.ErrSessionClosed
	}

	answer.SubmittedAt = time.Now()
	return r.repo.SaveAnswer(ctx, &answer)
}

// GetSessionResults получает результаты тестирования
func (r *TestRunner) GetSessionResults(ctx context.Context, sessionID string) (*models.TestResult, error) {
	result, err := r.repo.GetResult(ctx, sessionID)
	if err != nil {
		return nil, err
	}
	return result, nil
}

// GetActiveSessions возвращает активные сессии тестирования
func (r *TestRunner) GetActiveSessions(ctx context.Context) ([]models.TestSession, error) {
	return r.repo.ListActiveSessions(ctx)
}

// CleanupExpiredSessions очищает истекшие сессии
func (r *TestRunner) CleanupExpiredSessions(ctx context.Context) error {
	expirationTime := time.Now().Add(-24 * time.Hour)
	return r.repo.CleanupExpiredSessions(ctx, expirationTime)
}

// ExtendSessionTime продлевает время сессии
func (r *TestRunner) ExtendSessionTime(ctx context.Context, sessionID string, duration time.Duration) error {
	session, err := r.GetSession(ctx, sessionID)
	if err != nil {
		return err
	}

	if session.Status != "in_progress" {
		return errors.ErrSessionClosed
	}

	session.TimeoutAt = session.TimeoutAt.Add(duration)
	return r.repo.UpdateSession(ctx, session)
}
