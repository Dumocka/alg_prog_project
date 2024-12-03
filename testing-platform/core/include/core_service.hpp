#pragma once

#include <string>
#include <memory>
#include <vector>
#include <optional>
#include "models.hpp"
#include "auth/include/authorization_server.hpp"
#include "task_manager.hpp"
#include "test_runner.hpp"
#include "result_analyzer.hpp"

namespace core {

class CoreService {
public:
    CoreService(
        std::shared_ptr<auth::AuthorizationServer> auth_server,
        std::shared_ptr<TaskManager> task_manager,
        std::shared_ptr<TestRunner> test_runner,
        std::shared_ptr<ResultAnalyzer> result_analyzer
    );

    // Методы для работы с задачами
    std::vector<Task> get_available_tasks(const std::string& user_email);
    Task create_task(const TaskCreationRequest& request);
    void update_task(const Task& task);
    void delete_task(const std::string& task_id);

    // Методы для работы с тестами
    TestSession start_test_session(const std::string& user_email, const std::string& task_id);
    void submit_answer(const std::string& session_id, const Answer& answer);
    TestResult get_test_results(const std::string& session_id);
    
    // Методы для анализа результатов
    std::vector<TestResult> get_user_results(const std::string& user_email);
    AnalyticsReport generate_analytics(const AnalyticsRequest& request);
    
    // Методы для управления системой
    SystemStatus get_system_status();
    void update_system_settings(const SystemSettings& settings);

private:
    std::shared_ptr<auth::AuthorizationServer> auth_server_;
    std::shared_ptr<TaskManager> task_manager_;
    std::shared_ptr<TestRunner> test_runner_;
    std::shared_ptr<ResultAnalyzer> result_analyzer_;

    // Внутренние методы
    bool verify_user_access(const std::string& user_email, const std::string& task_id);
    void log_activity(const std::string& user_email, const std::string& action);
};

} // namespace core
