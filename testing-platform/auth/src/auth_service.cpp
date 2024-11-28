#include "auth_service.hpp"
#include <openssl/sha.h>
#include <openssl/hmac.h>
#include <jwt-cpp/jwt.h>
#include <ctime>
#include <sstream>
#include <iomanip>
#include <random>

namespace auth {

namespace {
    const std::string JWT_SECRET = "your-secret-key"; // In production, this should be loaded from environment
    const int TOKEN_EXPIRY = 3600; // 1 hour
    const int REFRESH_TOKEN_EXPIRY = 2592000; // 30 days
}

AuthService::AuthService(std::unique_ptr<Database> db)
    : db_(std::move(db)) {}

Token AuthService::authenticate(const std::string& email, const std::string& password) {
    auto user = db_->get_user_by_email(email);
    if (!user) {
        throw std::runtime_error("User not found");
    }

    if (hash_password(password) != user->password_hash) {
        throw std::runtime_error("Invalid password");
    }

    Token token;
    token.access_token = generate_token(*user);
    token.refresh_token = generate_refresh_token(*user);
    token.expires_in = TOKEN_EXPIRY;
    return token;
}

Token AuthService::refresh_token(const std::string& refresh_token) {
    try {
        auto decoded = jwt::decode(refresh_token);
        auto verifier = jwt::verify()
            .allow_algorithm(jwt::algorithm::hs256{JWT_SECRET});
        verifier.verify(decoded);

        auto user_id = decoded.get_payload_claim("user_id").as_int();
        auto user = db_->get_user_by_id(user_id);
        if (!user) {
            throw std::runtime_error("User not found");
        }

        Token token;
        token.access_token = generate_token(*user);
        token.refresh_token = generate_refresh_token(*user);
        token.expires_in = TOKEN_EXPIRY;
        return token;
    } catch (const std::exception& e) {
        throw std::runtime_error("Invalid refresh token");
    }
}

bool AuthService::verify_token(const std::string& token) {
    try {
        auto decoded = jwt::decode(token);
        auto verifier = jwt::verify()
            .allow_algorithm(jwt::algorithm::hs256{JWT_SECRET});
        verifier.verify(decoded);
        return true;
    } catch (const std::exception&) {
        return false;
    }
}

User AuthService::create_user(const std::string& email, const std::string& name, const std::string& password) {
    return db_->create_user(email, name, hash_password(password));
}

std::optional<User> AuthService::get_user(int64_t id) {
    return db_->get_user_by_id(id);
}

bool AuthService::check_permission(int64_t user_id, const std::string& resource, const std::string& action) {
    return db_->check_permission(user_id, resource, action);
}

void AuthService::assign_role(int64_t user_id, int64_t role_id) {
    db_->assign_role_to_user(user_id, role_id);
}

std::vector<Role> AuthService::get_user_roles(int64_t user_id) {
    return db_->get_user_roles(user_id);
}

std::string AuthService::hash_password(const std::string& password) {
    unsigned char hash[SHA256_DIGEST_LENGTH];
    SHA256_CTX sha256;
    SHA256_Init(&sha256);
    SHA256_Update(&sha256, password.c_str(), password.length());
    SHA256_Final(hash, &sha256);

    std::stringstream ss;
    for (int i = 0; i < SHA256_DIGEST_LENGTH; i++) {
        ss << std::hex << std::setw(2) << std::setfill('0') << static_cast<int>(hash[i]);
    }
    return ss.str();
}

std::string AuthService::generate_token(const User& user) {
    auto token = jwt::create()
        .set_issuer("auth-service")
        .set_type("JWS")
        .set_issued_at(std::chrono::system_clock::now())
        .set_expires_at(std::chrono::system_clock::now() + std::chrono::seconds{TOKEN_EXPIRY})
        .set_payload_claim("user_id", jwt::claim(std::to_string(user.id)))
        .set_payload_claim("email", jwt::claim(user.email))
        .sign(jwt::algorithm::hs256{JWT_SECRET});
    return token;
}

std::string AuthService::generate_refresh_token(const User& user) {
    auto token = jwt::create()
        .set_issuer("auth-service")
        .set_type("JWS")
        .set_issued_at(std::chrono::system_clock::now())
        .set_expires_at(std::chrono::system_clock::now() + std::chrono::seconds{REFRESH_TOKEN_EXPIRY})
        .set_payload_claim("user_id", jwt::claim(std::to_string(user.id)))
        .sign(jwt::algorithm::hs256{JWT_SECRET});
    return token;
}

} // namespace auth
