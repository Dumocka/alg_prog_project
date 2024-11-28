#include "database.hpp"
#include <chrono>
#include <ctime>
#include <sstream>

namespace auth {

Database::Database(const std::string& connection_string)
    : conn_(std::make_unique<pqxx::connection>(connection_string)) {}

std::optional<User> Database::get_user_by_email(const std::string& email) {
    try {
        pqxx::work txn{*conn_};
        auto result = txn.exec_params(
            "SELECT id, email, name, password_hash, created_at FROM users WHERE email = $1",
            email
        );
        txn.commit();

        if (result.empty()) {
            return std::nullopt;
        }

        User user;
        user.id = result[0][0].as<int64_t>();
        user.email = result[0][1].as<std::string>();
        user.name = result[0][2].as<std::string>();
        user.password_hash = result[0][3].as<std::string>();
        user.created_at = result[0][4].as<std::string>();
        return user;
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

std::optional<User> Database::get_user_by_id(int64_t id) {
    try {
        pqxx::work txn{*conn_};
        auto result = txn.exec_params(
            "SELECT id, email, name, password_hash, created_at FROM users WHERE id = $1",
            id
        );
        txn.commit();

        if (result.empty()) {
            return std::nullopt;
        }

        User user;
        user.id = result[0][0].as<int64_t>();
        user.email = result[0][1].as<std::string>();
        user.name = result[0][2].as<std::string>();
        user.password_hash = result[0][3].as<std::string>();
        user.created_at = result[0][4].as<std::string>();
        return user;
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

User Database::create_user(const std::string& email, const std::string& name, const std::string& password_hash) {
    try {
        pqxx::work txn{*conn_};
        auto result = txn.exec_params(
            "INSERT INTO users (email, name, password_hash, created_at) VALUES ($1, $2, $3, NOW()) RETURNING id, email, name, password_hash, created_at",
            email,
            name,
            password_hash
        );
        txn.commit();

        User user;
        user.id = result[0][0].as<int64_t>();
        user.email = result[0][1].as<std::string>();
        user.name = result[0][2].as<std::string>();
        user.password_hash = result[0][3].as<std::string>();
        user.created_at = result[0][4].as<std::string>();
        return user;
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

std::vector<Role> Database::get_user_roles(int64_t user_id) {
    try {
        pqxx::work txn{*conn_};
        auto result = txn.exec_params(
            "SELECT r.id, r.name, r.permissions FROM roles r "
            "JOIN user_roles ur ON ur.role_id = r.id "
            "WHERE ur.user_id = $1",
            user_id
        );
        txn.commit();

        std::vector<Role> roles;
        for (const auto& row : result) {
            Role role;
            role.id = row[0].as<int64_t>();
            role.name = row[1].as<std::string>();
            role.permissions = row[2].as<std::vector<std::string>>();
            roles.push_back(role);
        }
        return roles;
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

void Database::assign_role_to_user(int64_t user_id, int64_t role_id) {
    try {
        pqxx::work txn{*conn_};
        txn.exec_params(
            "INSERT INTO user_roles (user_id, role_id) VALUES ($1, $2) ON CONFLICT DO NOTHING",
            user_id,
            role_id
        );
        txn.commit();
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

bool Database::check_permission(int64_t user_id, const std::string& resource, const std::string& action) {
    try {
        pqxx::work txn{*conn_};
        auto result = txn.exec_params(
            "SELECT EXISTS ("
            "  SELECT 1 FROM user_roles ur "
            "  JOIN roles r ON r.id = ur.role_id "
            "  WHERE ur.user_id = $1 "
            "    AND r.permissions ? ($2 || ':' || $3)"
            ")",
            user_id,
            resource,
            action
        );
        txn.commit();
        return result[0][0].as<bool>();
    } catch (const std::exception& e) {
        throw std::runtime_error("Database error: " + std::string(e.what()));
    }
}

} // namespace auth
