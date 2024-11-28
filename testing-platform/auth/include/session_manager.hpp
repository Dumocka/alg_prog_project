#pragma once

#include <string>
#include <memory>
#include <cpp_redis/cpp_redis>
#include "models.hpp"

namespace auth {

class SessionManager {
public:
    explicit SessionManager(const std::string& redis_host = "127.0.0.1", int redis_port = 6379);
    ~SessionManager() = default;

    // Session management
    std::string create_session(const std::string& login_token);
    bool validate_session(const std::string& session_token);
    void update_session(const std::string& session_token, const User& user);
    void remove_session(const std::string& session_token);

    // Login token management
    std::string generate_login_token();
    bool validate_login_token(const std::string& login_token);

private:
    std::unique_ptr<cpp_redis::client> redis_client_;
    static constexpr int SESSION_TTL = 3600;  // 1 hour
    static constexpr int LOGIN_TOKEN_TTL = 300;  // 5 minutes
};

} // namespace auth
