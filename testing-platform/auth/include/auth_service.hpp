#pragma once

#include <string>
#include <memory>
#include "database.hpp"
#include "models.hpp"

namespace auth {

class AuthService {
public:
    explicit AuthService(std::unique_ptr<Database> db);
    ~AuthService() = default;

    // Authentication
    Token authenticate(const std::string& email, const std::string& password);
    Token refresh_token(const std::string& refresh_token);
    bool verify_token(const std::string& token);
    
    // User management
    User create_user(const std::string& email, const std::string& name, const std::string& password);
    std::optional<User> get_user(int64_t id);
    
    // Permission management
    bool check_permission(int64_t user_id, const std::string& resource, const std::string& action);
    void assign_role(int64_t user_id, int64_t role_id);
    std::vector<Role> get_user_roles(int64_t user_id);

private:
    std::string hash_password(const std::string& password);
    std::string generate_token(const User& user);
    std::string generate_refresh_token(const User& user);
    
    std::unique_ptr<Database> db_;
};

} // namespace auth
