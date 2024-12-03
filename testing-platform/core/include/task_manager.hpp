#pragma once

#include <string>
#include <vector>
#include <memory>
#include <optional>
#include "models.hpp"
#include "database.hpp"

namespace core {

class TaskManager {
public:
    explicit TaskManager(std::shared_ptr<Database> db);

    // Управление задачами
    std::vector<Task> get_tasks(const std::optional<std::string>& author_email = std::nullopt);
    std::optional<Task> get_task_by_id(const std::string& task_id);
    Task create_task(const TaskCreationRequest& request, const std::string& author_email);
    void update_task(const Task& task);
    void delete_task(const std::string& task_id);

    // Поиск и фильтрация
    std::vector<Task> search_tasks(const std::string& query);
    std::vector<Task> filter_tasks(const std::vector<std::string>& tags,
                                 const std::optional<std::string>& difficulty = std::nullopt);

    // Управление тегами
    std::vector<std::string> get_all_tags();
    void add_tag_to_task(const std::string& task_id, const std::string& tag);
    void remove_tag_from_task(const std::string& task_id, const std::string& tag);

private:
    std::shared_ptr<Database> db_;

    // Внутренние методы
    bool validate_task(const Task& task);
    void index_task(const Task& task);
};

} // namespace core
