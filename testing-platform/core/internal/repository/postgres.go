package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	_ "github.com/lib/pq"
	"testing-platform/core/internal/model"
)

type PostgresRepository struct {
	db *sql.DB
}

func NewPostgresRepository(url string) (*PostgresRepository, error) {
	db, err := sql.Open("postgres", url)
	if err != nil {
		return nil, err
	}
	
	if err := db.Ping(); err != nil {
		return nil, err
	}

	return &PostgresRepository{db: db}, nil
}

func (r *PostgresRepository) CreateTest(ctx context.Context, test *model.Test) error {
	query := `
		INSERT INTO tests (title, description, creator_id, is_public, created_at, updated_at)
		VALUES ($1, $2, $3, $4, $5, $6)
		RETURNING id`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		test.Title,
		test.Description,
		test.CreatorID,
		test.IsPublic,
		test.CreatedAt,
		test.UpdatedAt,
	).Scan(&test.ID)
}

func (r *PostgresRepository) GetTest(ctx context.Context, id int64) (*model.Test, error) {
	test := &model.Test{}
	query := `
		SELECT id, title, description, creator_id, is_public, created_at, updated_at
		FROM tests
		WHERE id = $1`
	
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&test.ID,
		&test.Title,
		&test.Description,
		&test.CreatorID,
		&test.IsPublic,
		&test.CreatedAt,
		&test.UpdatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return test, err
}

func (r *PostgresRepository) CreateQuestion(ctx context.Context, question *model.Question) error {
	options, err := json.Marshal(question.Options)
	if err != nil {
		return err
	}

	query := `
		INSERT INTO questions (test_id, type, text, options, required)
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		question.TestID,
		question.Type,
		question.Text,
		options,
		question.Required,
	).Scan(&question.ID)
}

func (r *PostgresRepository) SaveAnswer(ctx context.Context, answer *model.Answer) error {
	query := `
		INSERT INTO answers (question_id, user_id, value, created_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id`
	
	return r.db.QueryRowContext(
		ctx,
		query,
		answer.QuestionID,
		answer.UserID,
		answer.Value,
		answer.CreatedAt,
	).Scan(&answer.ID)
}

func (r *PostgresRepository) GetTestResults(ctx context.Context, testID, userID int64) (*model.TestResult, error) {
	result := &model.TestResult{}
	query := `
		SELECT id, test_id, user_id, score, created_at
		FROM test_results
		WHERE test_id = $1 AND user_id = $2`
	
	err := r.db.QueryRowContext(ctx, query, testID, userID).Scan(
		&result.ID,
		&result.TestID,
		&result.UserID,
		&result.Score,
		&result.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	return result, err
}

func (r *PostgresRepository) CheckTestPermission(ctx context.Context, testID, userID int64) (string, error) {
	var role string
	query := `
		SELECT role
		FROM test_permissions
		WHERE test_id = $1 AND user_id = $2`
	
	err := r.db.QueryRowContext(ctx, query, testID, userID).Scan(&role)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return role, err
}
