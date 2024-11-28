// Интерфейс для работы с базой данных
#pragma once

#include <memory>
#include <pqxx/pqxx>
#include "models.hpp"

namespace auth {

// Класс для работы с базой данных
class Database {
public:
    // Конструктор принимает строку подключения к базе данных
    explicit Database(const std::string& connection_string);
    ~Database() = default;

    // Операции с пользователями
    std::optional<User> get_user_by_email(const std::string& email);    // Получение пользователя по email
    std::optional<User> get_user_by_id(int64_t id);                     // Получение пользователя по ID
    User create_user(const std::string& email,                          // Создание нового пользователя
                    const std::string& name, 
                    const std::string& password_hash);
    
    // Операции с ролями
    std::vector<Role> get_user_roles(int64_t user_id);                 // Получение ролей пользователя
    void assign_role_to_user(int64_t user_id, int64_t role_id);        // Назначение роли пользователю
    
    // Операции с разрешениями
    bool check_permission(int64_t user_id,                             // Проверка разрешений пользователя
                         const std::string& resource, 
                         const std::string& action);

private:
    std::unique_ptr<pqxx::connection> conn_;    // Соединение с базой данных
};

} // namespace auth
