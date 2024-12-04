package service

import (
	"context"
	"time"

	"testing-platform/core/internal/model"
	"testing-platform/core/internal/repository"
	"testing-platform/core/internal/errors"
)

var (
	ErrTestNotFound   = errors.ErrTaskNotFound
	ErrInvalidAnswer  = errors.ErrInvalidAnswer
)

// TestService сервис для работы с тестами
type TestService struct {
	testRepo repository.TestRepository
}

// NewTestService создает новый экземпляр сервиса
func NewTestService(testRepo repository.TestRepository) *TestService {
	return &TestService{
		testRepo: testRepo,
	}
}

// CreateTest создает новый тест
func (s *TestService) CreateTest(ctx context.Context, userID int64, test *model.Test) error {
	test.CreatorID = userID
	test.CreatedAt = time.Now()
	test.UpdatedAt = time.Now()

	return s.testRepo.CreateTest(ctx, test)
}

// GetTest возвращает тест по ID
func (s *TestService) GetTest(ctx context.Context, userID, testID int64) (*model.Test, error) {
	test, err := s.testRepo.GetTestByID(ctx, testID)
	if err != nil {
		return nil, ErrTestNotFound
	}

	// Проверка прав доступа 
	if !test.IsPublic && test.CreatorID != userID {
		return nil, errors.ErrUnauthorized
	}

	return test, nil
}

// ListTests возвращает список тестов с фильтрацией
func (s *TestService) ListTests(ctx context.Context, userID int64, filters map[string]interface{}) ([]model.Test, error) {
	filters["creator_id"] = userID
	return s.testRepo.ListTests(ctx, filters)
}

// AddQuestion добавляет вопрос к тесту
func (s *TestService) AddQuestion(ctx context.Context, userID int64, question *model.Question) error {
	// Проверяем, что тест принадлежит пользователю
	test, err := s.testRepo.GetTestByID(ctx, question.TestID)
	if err != nil {
		return ErrTestNotFound
	}

	if test.CreatorID != userID {
		return errors.ErrUnauthorized
	}

	return s.testRepo.AddQuestion(ctx, question)
}

// SubmitAnswer обрабатывает ответ на вопрос
func (s *TestService) SubmitAnswer(ctx context.Context, userID int64, answer *model.Answer) error {
	// Проверяем корректность ответа
	if answer.Content == "" {
		return ErrInvalidAnswer
	}

	// В реальном приложении здесь будет логика проверки правильности ответа
	answer.IsCorrect = s.checkAnswer(answer)
	answer.UserID = userID

	return s.testRepo.SubmitAnswer(ctx, answer)
}

// checkAnswer проверяет правильность ответа (заглушка)
func (s *TestService) checkAnswer(answer *model.Answer) bool {
	// TODO: Реализовать проверку ответа
	return true
}

// GetResults возвращает результаты теста
func (s *TestService) GetResults(ctx context.Context, testID, userID int64) (*model.TestResult, error) {
	return s.testRepo.GetTestResults(ctx, testID, userID)
}
