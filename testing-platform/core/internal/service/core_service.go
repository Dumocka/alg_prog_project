package service

import (
	"context"

	"github.com/your-org/testing-platform/core/internal/errors"
	"github.com/your-org/testing-platform/core/internal/models"
	"github.com/your-org/testing-platform/core/internal/repository"
)

// CoreService представляет основной сервис платформы тестирования
type CoreService struct {
	repo           repository.Repository
	taskManager    *TaskManager
	testRunner     *TestRunner
	resultAnalyzer *ResultAnalyzer
}

// NewCoreService создает новый экземпляр CoreService
func NewCoreService(repo repository.Repository) *CoreService {
	return &CoreService{
		repo:           repo,
		taskManager:    NewTaskManager(repo),
		testRunner:     NewTestRunner(repo),
		resultAnalyzer: NewResultAnalyzer(repo),
	}
}

// GetAvailableTasks возвращает доступные задачи для пользователя
func (s *CoreService) GetAvailableTasks(ctx context.Context, userEmail string) ([]models.Task, error) {
	if userEmail == "" {
		return nil, errors.ErrInvalidUserEmail
	}

	filter := repository.TaskFilter{
		IsPublic: &[]bool{true}[0],
	}
	return s.taskManager.GetTasks(ctx, filter)
}

// CreateTask создает новую задачу
func (s *CoreService) CreateTask(ctx context.Context, req models.TaskCreationRequest, authorEmail string) (*models.Task, error) {
	return s.taskManager.CreateTask(ctx, req, authorEmail)
}

// StartTestSession начинает новую сессию тестирования
func (s *CoreService) StartTestSession(ctx context.Context, userEmail, taskID string) (*models.TestSession, error) {
	// Проверяем существование задачи
	task, err := s.taskManager.GetTaskByID(ctx, taskID)
	if err != nil {
		return nil, err
	}

	// Проверяем доступность задачи
	if !task.IsPublic {
		return nil, errors.ErrUnauthorized
	}

	return s.testRunner.CreateSession(ctx, userEmail, taskID)
}

// SubmitAnswer отправляет ответ на задание
func (s *CoreService) SubmitAnswer(ctx context.Context, answer models.Answer) error {
	if err := s.validateAnswer(ctx, answer); err != nil {
		return err
	}
	return s.testRunner.SubmitAnswer(ctx, answer)
}

// GetTestResults получает результаты тестирования
func (s *CoreService) GetTestResults(ctx context.Context, sessionID string) (*models.TestResult, error) {
	if sessionID == "" {
		return nil, errors.ErrInvalidRequest
	}
	return s.testRunner.GetSessionResults(ctx, sessionID)
}

// GenerateAnalytics генерирует аналитический отчет
func (s *CoreService) GenerateAnalytics(ctx context.Context, req models.AnalyticsRequest) (*models.AnalyticsReport, error) {
	return s.resultAnalyzer.GenerateReport(ctx, req)
}

// GetSystemStatus возвращает текущий статус системы
func (s *CoreService) GetSystemStatus(ctx context.Context) (*models.SystemStatus, error) {
	sessions, err := s.testRunner.GetActiveSessions(ctx)
	if err != nil {
		return nil, err
	}

	return &models.SystemStatus{
		IsHealthy:      true,
		Version:        "1.0.0",
		ActiveSessions: len(sessions),
		ActiveServices: []string{"core", "auth"},
	}, nil
}

// Внутренние методы

func (s *CoreService) validateAnswer(ctx context.Context, answer models.Answer) error {
	if answer.SessionID == "" || answer.Content == "" {
		return errors.ErrInvalidAnswer
	}

	session, err := s.testRunner.GetSession(ctx, answer.SessionID)
	if err != nil {
		return err
	}

	if session.Status != "in_progress" {
		return errors.ErrSessionClosed
	}

	return nil
}
