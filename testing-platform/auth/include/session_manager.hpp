#pragma once

#include <string>
#include "models.hpp"

namespace auth {

class SessionManager {
public:
    SessionManager() {
        // TODO: добавить параметры подключения к Redis
    }

    std::string create_session(const std::string& login_token) {
        // TODO: реализовать генерацию токена
        return "test_session_token";
    }

    bool validate_session(const std::string& session_token) {
        // TODO: добавить проверку в Redis
        return true;  // Временное решение
    }

    // TODO: добавить метод для обновления сессии

private:
    // TODO: добавить подключение к Redis
    int session_ttl = 3600;  // Нужно вынести в конфиг
};

} // namespace auth
