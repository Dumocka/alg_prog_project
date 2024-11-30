#pragma once

#include <string>
#include <vector>
#include <chrono>
#include "models.hpp"

namespace auth {

class JWTManager {
public:
    explicit JWTManager(const std::string& secret_key);

    // Создание токена доступа
    std::string create_access_token(
        const std::string& email,
        const std::vector<std::string>& permissions,
        const std::chrono::seconds& expiry = std::chrono::minutes(1)
    );

    // Создание токена обновления
    std::string create_refresh_token(
        const std::string& email,
        const std::chrono::seconds& expiry = std::chrono::hours(24 * 7)
    );

    // Проверка токена доступа
    bool verify_access_token(const std::string& token);

    // Проверка токена обновления
    bool verify_refresh_token(const std::string& token);

    // Получение email из токена обновления
    std::string get_email_from_refresh_token(const std::string& token);

    // Получение разрешений из токена доступа
    std::vector<std::string> get_permissions_from_access_token(const std::string& token);

private:
    std::string secret_key_;
    
    // Внутренние методы для работы с JWT
    std::string create_token(const nlohmann::json& payload);
    bool verify_token(const std::string& token);
    nlohmann::json decode_token(const std::string& token);
};

} // namespace auth
