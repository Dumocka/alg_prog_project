package service

import (
	"context"
	"time"

	"github.com/your-org/testing-platform/core/internal/errors"
	"github.com/your-org/testing-platform/core/internal/models"
	"github.com/your-org/testing-platform/core/internal/repository"
)

// ResultAnalyzer анализирует результаты тестирования
type ResultAnalyzer struct {
	repo repository.Repository
}

// NewResultAnalyzer создает новый анализатор результатов
func NewResultAnalyzer(repo repository.Repository) *ResultAnalyzer {
	return &ResultAnalyzer{repo: repo}
}

// GenerateReport генерирует аналитический отчет
func (a *ResultAnalyzer) GenerateReport(ctx context.Context, req models.AnalyticsRequest) (*models.AnalyticsReport, error) {
	if err := a.validateRequest(req); err != nil {
		return nil, err
	}

	filter := repository.ResultFilter{
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
	}

	if req.UserEmail != nil {
		filter.UserEmail = *req.UserEmail
	}
	if req.TaskID != nil {
		filter.TaskID = *req.TaskID
	}

	results, err := a.repo.GetResults(ctx, filter)
	if err != nil {
		return nil, errors.ErrAnalyticsFailure
	}

	metrics := make(map[string]interface{})
	for _, metric := range req.Metrics {
		switch metric {
		case "average_score":
			score, err := a.calculateAverageScore(results)
			if err != nil {
				return nil, err
			}
			metrics["average_score"] = score

		case "completion_rate":
			rate, err := a.calculateCompletionRate(results)
			if err != nil {
				return nil, err
			}
			metrics["completion_rate"] = rate

		case "difficulty_distribution":
			dist, err := a.calculateDifficultyDistribution(ctx, filter)
			if err != nil {
				return nil, err
			}
			metrics["difficulty_distribution"] = dist

		case "time_statistics":
			stats, err := a.calculateTimeStatistics(results)
			if err != nil {
				return nil, err
			}
			metrics["time_statistics"] = stats
		}
	}

	report := &models.AnalyticsReport{
		Metrics:     metrics,
		GeneratedAt: time.Now(),
	}

	if err := a.repo.SaveAnalyticsReport(ctx, report); err != nil {
		return nil, err
	}

	return report, nil
}

func (a *ResultAnalyzer) validateRequest(req models.AnalyticsRequest) error {
	if len(req.Metrics) == 0 {
		return errors.ErrInvalidRequest
	}

	if req.StartDate != nil && req.EndDate != nil {
		if req.StartDate.After(*req.EndDate) {
			return errors.ErrInvalidRequest
		}
	}

	return nil
}

func (a *ResultAnalyzer) calculateAverageScore(results []models.TestResult) (float64, error) {
	if len(results) == 0 {
		return 0, nil
	}

	var totalScore float64
	for _, result := range results {
		totalScore += result.Score
	}

	return totalScore / float64(len(results)), nil
}

func (a *ResultAnalyzer) calculateCompletionRate(results []models.TestResult) (float64, error) {
	if len(results) == 0 {
		return 0, nil
	}

	completed := 0
	for _, result := range results {
		if result.Status == "completed" {
			completed++
		}
	}

	return float64(completed) / float64(len(results)), nil
}

func (a *ResultAnalyzer) calculateDifficultyDistribution(ctx context.Context, filter repository.ResultFilter) (map[string]int, error) {
	tasks, err := a.repo.ListTasks(ctx, repository.TaskFilter{})
	if err != nil {
		return nil, err
	}

	distribution := make(map[string]int)
	for _, task := range tasks {
		distribution[task.Difficulty]++
	}

	return distribution, nil
}

func (a *ResultAnalyzer) calculateTimeStatistics(results []models.TestResult) (map[string]float64, error) {
	if len(results) == 0 {
		return map[string]float64{
			"average": 0,
			"minimum": 0,
			"maximum": 0,
		}, nil
	}

	var totalDuration float64
	minDuration := float64(results[0].CompletedAt.Sub(results[0].StartedAt).Seconds())
	maxDuration := minDuration

	for _, result := range results {
		duration := float64(result.CompletedAt.Sub(result.StartedAt).Seconds())
		totalDuration += duration

		if duration < minDuration {
			minDuration = duration
		}
		if duration > maxDuration {
			maxDuration = duration
		}
	}

	return map[string]float64{
		"average": totalDuration / float64(len(results)),
		"minimum": minDuration,
		"maximum": maxDuration,
	}, nil
}
