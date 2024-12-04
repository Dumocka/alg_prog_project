package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"testing-platform/core/internal/model"

	"github.com/lib/pq"
)

// TestRepository интерфейс для работы с тестами
type TestRepository interface {
	CreateTest(ctx context.Context, test *model.Test) error
	GetTestByID(ctx context.Context, testID int64) (*model.Test, error)
	ListTests(ctx context.Context, filters map[string]interface{}) ([]model.Test, error)
	UpdateTest(ctx context.Context, test *model.Test) error
	DeleteTest(ctx context.Context, testID int64) error

	AddQuestion(ctx context.Context, question *model.Question) error
	GetQuestionsByTestID(ctx context.Context, testID int64) ([]model.Question, error)

	SubmitAnswer(ctx context.Context, answer *model.Answer) error
	GetTestResults(ctx context.Context, testID, userID int64) (*model.TestResult, error)
}

// PostgresTestRepository реализация репозитория для PostgreSQL
type PostgresTestRepository struct {
	db *sql.DB
}

// NewPostgresTestRepository создает новый экземпляр репозитория
func NewPostgresTestRepository(db *sql.DB) *PostgresTestRepository {
	return &PostgresTestRepository{db: db}
}

// CreateTest создает новый тест
func (r *PostgresTestRepository) CreateTest(ctx context.Context, test *model.Test) error {
	query := `
		INSERT INTO tests (title, description, creator_id, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	now := time.Now()
	err := r.db.QueryRowContext(ctx, query,
		test.Title,
		test.Description,
		test.CreatorID,
		test.IsPublic,
		now,
		now,
	).Scan(&test.ID)

	return err
}

// GetTestByID возвращает тест по его ID
func (r *PostgresTestRepository) GetTestByID(ctx context.Context, testID int64) (*model.Test, error) {
	query := `
		SELECT id, title, description, creator_id, is_public, created_at, updated_at
		FROM tests
		WHERE id = $1
	`
	test := &model.Test{}
	err := r.db.QueryRowContext(ctx, query, testID).Scan(
		&test.ID,
		&test.Title,
		&test.Description,
		&test.CreatorID,
		&test.IsPublic,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	return test, nil
}

// ListTests возвращает список тестов с возможностью фильтрации
func (r *PostgresTestRepository) ListTests(ctx context.Context, filters map[string]interface{}) ([]model.Test, error) {
	query := `
		SELECT id, title, description, creator_id, is_public, created_at, updated_at
		FROM tests
		WHERE 1=1
	`
	var args []interface{}
	argPos := 1

	if creatorID, ok := filters["creator_id"].(int64); ok {
		query += fmt.Sprintf(" AND creator_id = $%d", argPos)
		args = append(args, creatorID)
		argPos++
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tests []model.Test
	for rows.Next() {
		var test model.Test
		err := rows.Scan(
			&test.ID,
			&test.Title,
			&test.Description,
			&test.CreatorID,
			&test.IsPublic,
			&test.CreatedAt,
			&test.UpdatedAt,
		)
		if err != nil {
			return nil, err
		}
		tests = append(tests, test)
	}

	return tests, nil
}

// AddQuestion добавляет вопрос к тесту
func (r *PostgresTestRepository) AddQuestion(ctx context.Context, question *model.Question) error {
	query := `
		INSERT INTO questions (test_id, text, type, score, options, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		question.TestID,
		question.Text,
		question.Type,
		question.Score,
		pq.Array(question.Options),
		time.Now(),
	).Scan(&question.ID)

	return err
}

// SubmitAnswer сохраняет ответ пользователя
func (r *PostgresTestRepository) SubmitAnswer(ctx context.Context, answer *model.Answer) error {
	query := `
		INSERT INTO answers (user_id, question_id, test_id, response, is_correct, created_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id
	`
	err := r.db.QueryRowContext(ctx, query,
		answer.UserID,
		answer.QuestionID,
		answer.TestID,
		answer.Content,
		answer.IsCorrect,
		time.Now(),
	).Scan(&answer.ID)

	return err
}

// GetTestResults возвращает результаты теста для пользователя
func (r *PostgresTestRepository) GetTestResults(ctx context.Context, testID, userID int64) (*model.TestResult, error) {
	query := `
		SELECT 
			tr.id, 
			tr.test_id, 
			tr.user_id, 
			tr.total_score, 
			tr.max_score, 
			tr.percentage, 
			tr.passed_at, 
			tr.status,
			t.title AS test_title
		FROM test_results tr
		JOIN tests t ON tr.test_id = t.id
		WHERE tr.test_id = $1 AND tr.user_id = $2
	`
	result := &model.TestResult{}
	err := r.db.QueryRowContext(ctx, query, testID, userID).Scan(
		&result.ID,
		&result.TestID,
		&result.UserID,
		&result.TotalScore,
		&result.MaxScore,
		&result.Percentage,
		&result.PassedAt,
		&result.Status,
		&result.TestTitle,
	)
	if err != nil {
		return nil, err
	}

	return result, nil
}
