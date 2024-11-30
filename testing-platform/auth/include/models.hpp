// Модели данных для сервиса аутентификации
#pragma once

#include <string>
#include <vector>
#include <optional>
#include <nlohmann/json.hpp>

namespace auth {

// Структура пользователя
struct User {
    int64_t id;                  // Уникальный идентификатор
    std::string email;           // Электронная почта
    std::string name;            // Имя пользователя
    std::string password_hash;   // Хеш пароля
    std::string created_at;      // Дата создания
    std::vector<std::string> roles;          // Роли пользователя
    std::vector<std::string> refresh_tokens; // Токены обновления

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(User, id, email, name, password_hash, created_at, roles, refresh_tokens)
};

// Структура разрешения
struct Permission {
    std::string resource;     // Например: "user", "course"
    std::string action;       // Например: "read", "write"
    std::string scope;        // Например: "own", "all"

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Permission, resource, action, scope)
};

// Структура роли
struct Role {
    std::string name;                   // Название роли
    std::vector<Permission> permissions; // Список разрешений

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(Role, name, permissions)
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

// Структура для JWT payload
struct JWTPayload {
    std::string email;
    std::vector<std::string> permissions;
    int64_t exp;  // Время истечения токена

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(JWTPayload, email, permissions, exp)
};

// Структура для ответа OAuth провайдера
struct OAuthUserInfo {
    std::string email;
    std::string name;
    std::optional<std::string> avatar_url;

    NLOHMANN_DEFINE_TYPE_INTRUSIVE(OAuthUserInfo, email, name, avatar_url)
};

} // namespace auth
