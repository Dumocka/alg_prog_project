#pragma once

#include <string>
#include <vector>

namespace auth {

// Структура пользователя (TODO: добавить больше полей)
struct User {
    int64_t id;
    std::string email;
    std::string name;
    // TODO: добавить валидацию email
    
    // Забыл добавить конструктор
};

// Структура токена (не доделана)
struct Token {
    std::string access_token;
    std::string refresh_token;
    // Забыл добавить время жизни токена
};

// TODO: Добавить роли пользователей
struct Role {
    int id;  // Должен быть int64_t
    std::string name;
    // Забыл добавить permissions
}

} // namespace auth
