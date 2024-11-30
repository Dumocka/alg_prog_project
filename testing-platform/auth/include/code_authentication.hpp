#pragma once

#include <string>
#include <chrono>
#include <optional>
#include <memory>
#include <cpp_redis/cpp_redis>

namespace auth {

class CodeAuthentication {
public:
    explicit CodeAuthentication(std::shared_ptr<cpp_redis::client> redis_client);

    // Генерация кода авторизации
    std::string generate_code(const std::string& login_token);

    // Проверка кода
    bool verify_code(const std::string& code, const std::string& login_token);

    // Получение email пользователя по коду
    std::optional<std::string> get_email_by_code(const std::string& code);

    // Удаление использованного кода
    void invalidate_code(const std::string& code);

private:
    std::shared_ptr<cpp_redis::client> redis_client_;
    const std::chrono::seconds code_expiry_{300}; // 5 минут

    // Внутренние методы
    std::string generate_random_code();
    std::string get_redis_key(const std::string& code);
};

} // namespace auth
