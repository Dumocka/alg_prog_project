package service

import (
	"context"
	"errors"
	"time"
	"testing-platform/core/internal/model"
)

var (
	ErrUnauthorized = errors.New("unauthorized")
	ErrNotFound     = errors.New("not found")
)

type Repository interface {
	CreateTest(ctx context.Context, test *model.Test) error
	GetTest(ctx context.Context, id int64) (*model.Test, error)
	CreateQuestion(ctx context.Context, question *model.Question) error
	SaveAnswer(ctx context.Context, answer *model.Answer) error
	GetTestResults(ctx context.Context, testID, userID int64) (*model.TestResult, error)
	CheckTestPermission(ctx context.Context, testID, userID int64) (string, error)
}

type Service struct {
	repo Repository
}

func NewService(repo Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) CreateTest(ctx context.Context, userID int64, test *model.Test) error {
	test.CreatorID = userID
	test.CreatedAt = time.Now()
	test.UpdatedAt = test.CreatedAt
	return s.repo.CreateTest(ctx, test)
}

func (s *Service) GetTest(ctx context.Context, userID, testID int64) (*model.Test, error) {
	test, err := s.repo.GetTest(ctx, testID)
	if err != nil {
		return nil, err
	}
	if test == nil {
		return nil, ErrNotFound
	}

	if !test.IsPublic {
		role, err := s.repo.CheckTestPermission(ctx, testID, userID)
		if err != nil {
			return nil, err
		}
		if role == "" {
			return nil, ErrUnauthorized
		}
	}

	return test, nil
}

func (s *Service) AddQuestion(ctx context.Context, userID int64, question *model.Question) error {
	role, err := s.repo.CheckTestPermission(ctx, question.TestID, userID)
	if err != nil {
		return err
	}
	if role != "owner" && role != "editor" {
		return ErrUnauthorized
	}

	return s.repo.CreateQuestion(ctx, question)
}

func (s *Service) SubmitAnswer(ctx context.Context, userID int64, answer *model.Answer) error {
	answer.UserID = userID
	answer.CreatedAt = time.Now()
	return s.repo.SaveAnswer(ctx, answer)
}

func (s *Service) GetResults(ctx context.Context, userID, testID int64) (*model.TestResult, error) {
	role, err := s.repo.CheckTestPermission(ctx, testID, userID)
	if err != nil {
		return nil, err
	}
	if role == "" {
		return nil, ErrUnauthorized
	}

	return s.repo.GetTestResults(ctx, testID, userID)
}
