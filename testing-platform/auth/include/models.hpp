// Модели данных для сервиса аутентификации
#pragma once

#include <string>
#include <vector>
#include <nlohmann/json.hpp>

namespace auth {

// Структура пользователя
struct User {
    int64_t id;                  // Уникальный идентификатор
    std::string email;           // Электронная почта
    std::string name;            // Имя пользователя
    std::string password_hash;   // Хеш пароля
    std::string created_at;      // Дата создания

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(User, id, email, name, created_at)
};

// Структура роли пользователя
struct Role {
    int64_t id;                         // Уникальный идентификатор роли
    std::string name;                   // Название роли
    std::vector<std::string> permissions; // Список разрешений

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Role, id, name, permissions)
};

// Связь пользователя и роли
struct UserRole {
    int64_t user_id;    // ID пользователя
    int64_t role_id;    // ID роли

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(UserRole, user_id, role_id)
};

// Структура токена доступа
struct Token {
    std::string access_token;    // Токен доступа
    std::string refresh_token;   // Токен обновления
    int64_t expires_in;          // Время истечения в секундах

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Token, access_token, refresh_token, expires_in)
};

// Запрос на аутентификацию
struct AuthRequest {
    std::string email;     // Электронная почта
    std::string password;  // Пароль

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(AuthRequest, email, password)
};

// Проверка разрешений
struct PermissionCheck {
    int64_t user_id;        // ID пользователя
    std::string resource;   // Ресурс для проверки
    std::string action;     // Действие для проверки

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(PermissionCheck, user_id, resource, action)
};

} // namespace auth
