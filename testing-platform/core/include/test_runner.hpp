#pragma once

#include <string>
#include <memory>
#include <unordered_map>
#include <chrono>
#include "models.hpp"
#include "database.hpp"

namespace core {

class TestRunner {
public:
    TestRunner(std::shared_ptr<Database> db);

    // Управление сессиями тестирования
    TestSession create_session(const std::string& user_email, const std::string& task_id);
    TestSession get_session(const std::string& session_id);
    void update_session(const TestSession& session);
    void terminate_session(const std::string& session_id);

    // Обработка ответов
    void submit_answer(const Answer& answer);
    TestResult get_session_results(const std::string& session_id);
    
    // Управление временем
    void extend_session_time(const std::string& session_id, 
                           const std::chrono::minutes& additional_time);
    std::chrono::seconds get_remaining_time(const std::string& session_id);

    // Мониторинг
    std::vector<TestSession> get_active_sessions();
    void cleanup_expired_sessions();

private:
    std::shared_ptr<Database> db_;
    std::unordered_map<std::string, TestSession> active_sessions_;

    // Внутренние методы
    bool validate_answer(const Answer& answer);
    void grade_answer(const std::string& session_id, const Answer& answer);
    void check_session_timeout(const std::string& session_id);
};

} // namespace core
