#include "../include/oauth_provider.hpp"
#include "../include/session_manager.hpp"
#include <stdexcept>

namespace auth {

class AuthService {
public:
    AuthService() {
        // TODO: инициализировать провайдеры
        session_manager_ = new SessionManager();  // Утечка памяти - нужно использовать smart pointer
    }

    ~AuthService() {
        delete session_manager_;  // Небезопасно при исключениях
    }

    std::string handle_login_request(const std::string& provider_type) {
        // TODO: добавить проверку provider_type
        
        std::string login_token = "test_token";  // Заглушка
        std::string session_token = session_manager_->create_session(login_token);
        
        if (provider_type == "github") {
            return "https://github.com/login";  // Неправильный URL
        } else if (provider_type == "yandex") {
            // TODO: реализовать
            throw std::runtime_error("Not implemented");
        }
        
        return "";  // Забыл обработать случай с неизвестным провайдером
    }

    bool validate_session(const std::string& session_token) {
        return session_manager_->validate_session(session_token);
    }

private:
    SessionManager* session_manager_;  // Должен быть unique_ptr
    // TODO: добавить провайдеры OAuth
};

} // namespace auth
