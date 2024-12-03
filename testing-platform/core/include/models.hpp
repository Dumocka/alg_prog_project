#pragma once

#include <string>
#include <vector>
#include <chrono>
#include <optional>
#include <nlohmann/json.hpp>

namespace core {

// Модель задачи
struct Task {
    std::string id;
    std::string title;
    std::string description;
    std::string difficulty;
    std::vector<std::string> tags;
    std::string author_email;
    std::chrono::system_clock::time_point created_at;
    bool is_public;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Task, id, title, description, difficulty, tags, author_email, created_at, is_public)
};

// Запрос на создание задачи
struct TaskCreationRequest {
    std::string title;
    std::string description;
    std::string difficulty;
    std::vector<std::string> tags;
    bool is_public;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(TaskCreationRequest, title, description, difficulty, tags, is_public)
};

// Сессия тестирования
struct TestSession {
    std::string id;
    std::string user_email;
    std::string task_id;
    std::chrono::system_clock::time_point started_at;
    std::optional<std::chrono::system_clock::time_point> completed_at;
    std::string status; // "in_progress", "completed", "timeout"

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(TestSession, id, user_email, task_id, started_at, completed_at, status)
};

// Ответ пользователя
struct Answer {
    std::string session_id;
    std::string question_id;
    std::string content;
    std::chrono::system_clock::time_point submitted_at;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Answer, session_id, question_id, content, submitted_at)
};

// Результат тестирования
struct TestResult {
    std::string session_id;
    std::string user_email;
    std::string task_id;
    int score;
    std::vector<std::string> correct_answers;
    std::vector<std::string> wrong_answers;
    std::chrono::system_clock::time_point completed_at;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(TestResult, session_id, user_email, task_id, score, correct_answers, wrong_answers, completed_at)
};

// Запрос на аналитику
struct AnalyticsRequest {
    std::optional<std::string> user_email;
    std::optional<std::string> task_id;
    std::optional<std::chrono::system_clock::time_point> start_date;
    std::optional<std::chrono::system_clock::time_point> end_date;
    std::vector<std::string> metrics;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(AnalyticsRequest, user_email, task_id, start_date, end_date, metrics)
};

// Отчет по аналитике
struct AnalyticsReport {
    std::string report_id;
    nlohmann::json metrics;
    std::chrono::system_clock::time_point generated_at;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(AnalyticsReport, report_id, metrics, generated_at)
};

// Статус системы
struct SystemStatus {
    bool is_healthy;
    std::string version;
    int active_sessions;
    double cpu_usage;
    double memory_usage;
    std::vector<std::string> active_services;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(SystemStatus, is_healthy, version, active_sessions, cpu_usage, memory_usage, active_services)
};

// Настройки системы
struct SystemSettings {
    int max_concurrent_sessions;
    int session_timeout_minutes;
    bool maintenance_mode;
    std::vector<std::string> enabled_features;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(SystemSettings, max_concurrent_sessions, session_timeout_minutes, maintenance_mode, enabled_features)
};

} // namespace core
