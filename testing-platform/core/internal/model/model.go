package model

import "time"

// Test представляет тест
type Test struct {
    ID          int64     `json:"id" bson:"_id"`
    Title       string    `json:"title" bson:"title"`
    Description string    `json:"description" bson:"description"`
    CreatorID   int64     `json:"creator_id" bson:"creator_id"`
    IsPublic    bool      `json:"is_public" bson:"is_public"`
    CreatedAt   time.Time `json:"created_at" bson:"created_at"`
    UpdatedAt   time.Time `json:"updated_at" bson:"updated_at"`
}

// Question представляет вопрос в тесте
type Question struct {
    ID         int64    `json:"id" bson:"_id"`
    TestID     int64    `json:"test_id" bson:"test_id"`
    Text       string   `json:"text" bson:"text"`
    Type       string   `json:"type" bson:"type"` // multiple_choice, single_choice, text
    Options    []string `json:"options" bson:"options"`
    Difficulty string   `json:"difficulty" bson:"difficulty"`
    Score      int      `json:"score" bson:"score"`
}

// Answer представляет ответ пользователя
type Answer struct {
    ID          int64     `json:"id" bson:"_id"`
    SessionID   string    `json:"session_id" bson:"session_id"`
    TestID      int64     `json:"test_id" bson:"test_id"`
    QuestionID  int64     `json:"question_id" bson:"question_id"`
    UserID      int64     `json:"user_id" bson:"user_id"`
    Content     string    `json:"content" bson:"content"`
    IsCorrect   bool      `json:"is_correct" bson:"is_correct"`
    CreatedAt   time.Time `json:"created_at" bson:"created_at"`
}

// TestResult представляет результат тестирования
type TestResult struct {
    ID             int64     `json:"id" bson:"_id"`
    TestID         int64     `json:"test_id" bson:"test_id"`
    UserID         int64     `json:"user_id" bson:"user_id"`
    TestTitle      string    `json:"test_title" bson:"test_title"`
    TotalScore     int       `json:"total_score" bson:"total_score"`
    MaxScore       int       `json:"max_score" bson:"max_score"`
    Percentage     float64   `json:"percentage" bson:"percentage"`
    PassedAt       time.Time `json:"passed_at" bson:"passed_at"`
    Status         string    `json:"status" bson:"status"`
    TotalQuestions int       `json:"total_questions" bson:"total_questions"`
    CorrectAnswers int       `json:"correct_answers" bson:"correct_answers"`
    CompletedAt    time.Time `json:"completed_at" bson:"completed_at"`
}
