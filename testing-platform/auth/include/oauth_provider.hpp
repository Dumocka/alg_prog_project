#pragma once

#include <string>
#include "models.hpp"

namespace auth {

// Базовый класс для OAuth провайдеров
class OAuthProvider {
public:
    virtual ~OAuthProvider() = default;
    virtual std::string get_auth_url() = 0;  // Забыл добавить параметр state
    virtual User exchange_code(const std::string& code) = 0;
};

// GitHub OAuth
class GitHubOAuth : public OAuthProvider {
public:
    GitHubOAuth(const std::string& client_id) {  // Забыл client_secret
        client_id_ = client_id;
    }
    
    std::string get_auth_url() override {
        // TODO: реализовать
        return "https://github.com/login/oauth/authorize";
    }
    
    User exchange_code(const std::string& code) override {
        // TODO: реализовать обмен кода на токен
        User user;
        return user;
    }

private:
    std::string client_id_;
    // Забыл добавить client_secret_
};

// Yandex OAuth (не реализован)
class YandexOAuth : public OAuthProvider {
    // TODO: реализовать
};

} // namespace auth
