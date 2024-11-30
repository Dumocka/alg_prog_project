#pragma once

#include <string>
#include <memory>
#include <functional>
#include <utility>
#include "models.hpp"

namespace auth {

class OAuthProvider {
public:
    OAuthProvider(
        const std::string& client_id,
        const std::string& client_secret,
        const std::string& redirect_uri
    ) : client_id_(client_id), 
        client_secret_(client_secret),
        redirect_uri_(redirect_uri) {}

    virtual ~OAuthProvider() = default;
    
    // Получение URL для авторизации
    virtual std::string get_authorization_url(const std::string& state) = 0;
    
    // Обмен кода на токены
    virtual std::pair<std::string, std::string> exchange_code(const std::string& code) = 0;
    
    // Получение информации о пользователе
    virtual OAuthUserInfo get_user_info(const std::string& access_token) = 0;
    
    // Обновление токена доступа
    virtual std::pair<std::string, std::string> refresh_access_token(
        const std::string& refresh_token) = 0;

protected:
    std::string client_id_;
    std::string client_secret_;
    std::string redirect_uri_;
};

class GitHubOAuthProvider : public OAuthProvider {
public:
    GitHubOAuthProvider(
        const std::string& client_id,
        const std::string& client_secret,
        const std::string& redirect_uri
    );
    
    std::string get_authorization_url(const std::string& state) override;
    std::pair<std::string, std::string> exchange_code(const std::string& code) override;
    OAuthUserInfo get_user_info(const std::string& access_token) override;
    std::pair<std::string, std::string> refresh_access_token(
        const std::string& refresh_token) override;
};

class YandexOAuthProvider : public OAuthProvider {
public:
    YandexOAuthProvider(
        const std::string& client_id,
        const std::string& client_secret,
        const std::string& redirect_uri
    );
    
    std::string get_authorization_url(const std::string& state) override;
    std::pair<std::string, std::string> exchange_code(const std::string& code) override;
    OAuthUserInfo get_user_info(const std::string& access_token) override;
    std::pair<std::string, std::string> refresh_access_token(
        const std::string& refresh_token) override;
};

} // namespace auth
