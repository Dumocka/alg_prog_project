package errors

import "errors"

// Core service errors
var (
	ErrSessionClosed     = errors.New("test session is closed")
	ErrSessionNotFound   = errors.New("test session not found")
	ErrInvalidTaskID     = errors.New("invalid task ID")
	ErrInvalidUserEmail  = errors.New("invalid user email")
	ErrUnauthorized      = errors.New("unauthorized access")
	ErrTaskNotFound      = errors.New("task not found")
	ErrInvalidAnswer     = errors.New("invalid answer format")
	ErrDatabaseError     = errors.New("database operation failed")
	ErrInvalidRequest    = errors.New("invalid request parameters")
	ErrAnalyticsFailure  = errors.New("failed to generate analytics")
)
