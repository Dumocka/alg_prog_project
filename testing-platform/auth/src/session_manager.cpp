#include "session_manager.hpp"
#include <random>
#include <sstream>
#include <iomanip>
#include <nlohmann/json.hpp>

namespace auth {

namespace {
    std::string generate_random_token(size_t length) {
        static const std::string chars = 
            "0123456789"
            "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
            "abcdefghijklmnopqrstuvwxyz";
        
        std::random_device rd;
        std::mt19937 gen(rd());
        std::uniform_int_distribution<> dis(0, chars.length() - 1);
        
        std::string token;
        token.reserve(length);
        for (size_t i = 0; i < length; ++i) {
            token += chars[dis(gen)];
        }
        return token;
    }
}

SessionManager::SessionManager(const std::string& redis_host, int redis_port) {
    redis_client_ = std::make_unique<cpp_redis::client>();
    redis_client_->connect(redis_host, redis_port);
}

std::string SessionManager::create_session(const std::string& login_token) {
    std::string session_token = generate_random_token(32);
    
    nlohmann::json session_data = {
        {"status", "anonymous"},
        {"login_token", login_token}
    };
    
    redis_client_->set(session_token, session_data.dump(), [](cpp_redis::reply& reply) {
        if (!reply.is_string()) {
            throw std::runtime_error("Failed to create session");
        }
    });
    redis_client_->expire(session_token, SESSION_TTL);
    redis_client_->sync_commit();
    
    return session_token;
}

bool SessionManager::validate_session(const std::string& session_token) {
    bool valid = false;
    redis_client_->exists({session_token}, [&valid](cpp_redis::reply& reply) {
        valid = reply.as_integer() == 1;
    });
    redis_client_->sync_commit();
    return valid;
}

void SessionManager::update_session(const std::string& session_token, const User& user) {
    nlohmann::json session_data = {
        {"status", "authenticated"},
        {"user_id", user.id},
        {"email", user.email},
        {"name", user.name}
    };
    
    redis_client_->set(session_token, session_data.dump(), [](cpp_redis::reply& reply) {
        if (!reply.is_string()) {
            throw std::runtime_error("Failed to update session");
        }
    });
    redis_client_->expire(session_token, SESSION_TTL);
    redis_client_->sync_commit();
}

void SessionManager::remove_session(const std::string& session_token) {
    redis_client_->del({session_token}, [](cpp_redis::reply& reply) {
        if (!reply.is_integer() || reply.as_integer() != 1) {
            throw std::runtime_error("Failed to remove session");
        }
    });
    redis_client_->sync_commit();
}

std::string SessionManager::generate_login_token() {
    std::string login_token = generate_random_token(32);
    redis_client_->set(
        "login:" + login_token,
        "pending",
        [](cpp_redis::reply& reply) {
            if (!reply.is_string()) {
                throw std::runtime_error("Failed to create login token");
            }
        }
    );
    redis_client_->expire("login:" + login_token, LOGIN_TOKEN_TTL);
    redis_client_->sync_commit();
    return login_token;
}

bool SessionManager::validate_login_token(const std::string& login_token) {
    bool valid = false;
    redis_client_->get("login:" + login_token, [&valid](cpp_redis::reply& reply) {
        valid = reply.is_string() && reply.as_string() == "pending";
    });
    redis_client_->sync_commit();
    return valid;
}

} // namespace auth
