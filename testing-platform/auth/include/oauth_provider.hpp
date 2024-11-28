#pragma once

#include <string>
#include <memory>
#include <functional>
#include "models.hpp"

namespace auth {

class OAuthProvider {
public:
    virtual ~OAuthProvider() = default;
    
    // Get the authorization URL for the provider
    virtual std::string get_auth_url(const std::string& state) const = 0;
    
    // Exchange authorization code for user info
    virtual User exchange_code(const std::string& code) = 0;
};

class GitHubOAuth : public OAuthProvider {
public:
    GitHubOAuth(const std::string& client_id, const std::string& client_secret);
    
    std::string get_auth_url(const std::string& state) const override;
    User exchange_code(const std::string& code) override;
    
private:
    std::string client_id_;
    std::string client_secret_;
};

class YandexOAuth : public OAuthProvider {
public:
    YandexOAuth(const std::string& client_id, const std::string& client_secret);
    
    std::string get_auth_url(const std::string& state) const override;
    User exchange_code(const std::string& code) override;
    
private:
    std::string client_id_;
    std::string client_secret_;
};

} // namespace auth
