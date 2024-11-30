#pragma once

#include <string>
#include <memory>
#include <unordered_map>
#include <chrono>
#include "models.hpp"
#include "oauth_provider.hpp"
#include "code_authentication.hpp"
#include "jwt_manager.hpp"
#include "database.hpp"

namespace auth {

// Структура для хранения состояния авторизации
struct AuthState {
    std::chrono::system_clock::time_point expires_at;
    std::string status;  // "pending", "denied", "granted"
    std::optional<std::string> access_token;
    std::optional<std::string> refresh_token;
};

class AuthorizationServer {
public:
    AuthorizationServer(
        std::shared_ptr<Database> db,
        std::shared_ptr<JWTManager> jwt_manager,
        std::shared_ptr<CodeAuthentication> code_auth
    );

    // Инициация OAuth авторизации
    std::string initiate_oauth(const std::string& provider_type, const std::string& login_token);
    
    // Обработка OAuth callback
    void handle_oauth_callback(const std::string& provider_type, 
                             const std::string& code,
                             const std::string& state,
                             const std::string& error = "");

    // Инициация Code авторизации
    std::string initiate_code_auth(const std::string& login_token);
    
    // Проверка статуса авторизации
    AuthState check_auth_status(const std::string& login_token);
    
    // Обновление токена
    std::pair<std::string, std::string> refresh_tokens(const std::string& refresh_token);
    
    // Выход из системы
    void logout(const std::string& refresh_token);

private:
    std::shared_ptr<Database> db_;
    std::shared_ptr<JWTManager> jwt_manager_;
    std::shared_ptr<CodeAuthentication> code_auth_;
    std::unordered_map<std::string, AuthState> auth_states_;
    
    // OAuth провайдеры
    std::unique_ptr<OAuthProvider> github_provider_;
    std::unique_ptr<OAuthProvider> yandex_provider_;
    
    // Внутренние методы
    void create_or_update_user(const std::string& email, const std::string& name);
    std::pair<std::string, std::string> generate_tokens(const std::string& email);
    void cleanup_expired_states();
};

} // namespace auth
