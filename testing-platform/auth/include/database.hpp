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
    virtual ~Database() = default;

    // Операции с пользователями
    virtual std::optional<User> get_user_by_email(const std::string& email);    // Получение пользователя по email
    virtual std::optional<User> get_user_by_id(int64_t id);                     // Получение пользователя по ID
    virtual User create_user(const std::string& email,                          // Создание нового пользователя
                    const std::string& name, 
                    const std::string& password_hash);
    virtual void update_user(const User& user);                                 // Обновление пользователя
    
    // Операции с ролями
    virtual std::vector<Role> get_roles();                                      // Получение всех ролей
    virtual std::optional<Role> get_role_by_name(const std::string& name);      // Получение роли по имени
    virtual std::vector<Role> get_user_roles(int64_t user_id);                 // Получение ролей пользователя
    virtual void assign_role_to_user(const std::string& email,                 // Назначение роли пользователю
                                   const std::string& role_name);
    virtual void remove_role_from_user(const std::string& email,               // Удаление роли у пользователя
                                     const std::string& role_name);
    
    // Операции с разрешениями
    virtual std::vector<Permission> get_permissions_for_role(                  // Получение разрешений роли
        const std::string& role_name);
    virtual bool check_permission(int64_t user_id,                            // Проверка разрешений пользователя
                                const std::string& resource, 
                                const std::string& action);

    // Операции с refresh токенами
    virtual void add_refresh_token(const std::string& email,                  // Добавление refresh токена
                                 const std::string& token);
    virtual void remove_refresh_token(const std::string& email,               // Удаление refresh токена
                                    const std::string& token);
    virtual bool verify_refresh_token(const std::string& email,               // Проверка refresh токена
                                    const std::string& token);

protected:
    std::unique_ptr<pqxx::connection> conn_;    // Соединение с базой данных
};

// Фабричный метод для создания экземпляра базы данных
std::unique_ptr<Database> create_database(const std::string& connection_string);

} // namespace auth
