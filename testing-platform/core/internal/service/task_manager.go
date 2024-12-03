package service

import (
	"context"
	"time"

	"github.com/your-org/testing-platform/core/internal/errors"
	"github.com/your-org/testing-platform/core/internal/models"
	"github.com/your-org/testing-platform/core/internal/repository"
)

// TaskManager управляет задачами в системе
type TaskManager struct {
	repo repository.Repository
}

// NewTaskManager создает новый менеджер задач
func NewTaskManager(repo repository.Repository) *TaskManager {
	return &TaskManager{repo: repo}
}

// GetTasks возвращает список задач
func (m *TaskManager) GetTasks(ctx context.Context, filter repository.TaskFilter) ([]models.Task, error) {
	return m.repo.ListTasks(ctx, filter)
}

// CreateTask создает новую задачу
func (m *TaskManager) CreateTask(ctx context.Context, req models.TaskCreationRequest, authorEmail string) (*models.Task, error) {
	if authorEmail == "" {
		return nil, errors.ErrInvalidUserEmail
	}

	task := &models.Task{
		Title:       req.Title,
		Description: req.Description,
		Difficulty:  req.Difficulty,
		Tags:        req.Tags,
		AuthorEmail: authorEmail,
		CreatedAt:   time.Now(),
		IsPublic:    req.IsPublic,
	}

	if err := m.repo.CreateTask(ctx, task); err != nil {
		return nil, err
	}

	return task, nil
}

// GetTaskByID получает задачу по ID
func (m *TaskManager) GetTaskByID(ctx context.Context, taskID string) (*models.Task, error) {
	if taskID == "" {
		return nil, errors.ErrInvalidTaskID
	}

	task, err := m.repo.GetTask(ctx, taskID)
	if err != nil {
		return nil, errors.ErrTaskNotFound
	}

	return task, nil
}

// UpdateTask обновляет существующую задачу
func (m *TaskManager) UpdateTask(ctx context.Context, task models.Task) error {
	if task.ID == "" {
		return errors.ErrInvalidTaskID
	}

	return m.repo.UpdateTask(ctx, &task)
}

// DeleteTask удаляет задачу
func (m *TaskManager) DeleteTask(ctx context.Context, taskID string) error {
	if taskID == "" {
		return errors.ErrInvalidTaskID
	}

	return m.repo.DeleteTask(ctx, taskID)
}

// SearchTasks ищет задачи по запросу
func (m *TaskManager) SearchTasks(ctx context.Context, query string) ([]models.Task, error) {
	if query == "" {
		return nil, errors.ErrInvalidRequest
	}

	return m.repo.SearchTasks(ctx, query)
}

// GetAllTags возвращает все уникальные теги
func (m *TaskManager) GetAllTags(ctx context.Context) ([]string, error) {
	tasks, err := m.repo.ListTasks(ctx, repository.TaskFilter{})
	if err != nil {
		return nil, err
	}

	// Используем map для хранения уникальных тегов
	tagsMap := make(map[string]struct{})
	for _, task := range tasks {
		for _, tag := range task.Tags {
			tagsMap[tag] = struct{}{}
		}
	}

	// Преобразуем map в slice
	tags := make([]string, 0, len(tagsMap))
	for tag := range tagsMap {
		tags = append(tags, tag)
	}

	return tags, nil
}
