#include <cstdlib>
#include <memory>
#include <restbed>
#include <nlohmann/json.hpp>
#include "auth_service.hpp"
#include "database.hpp"

using json = nlohmann::json;
using namespace restbed;
using namespace auth;

namespace {
    std::unique_ptr<AuthService> auth_service;

    void authenticate_handler(const std::shared_ptr<Session> session) {
        const auto request = session->get_request();
        size_t content_length = request->get_header("Content-Length", 0);

        session->fetch(content_length, [](const std::shared_ptr<Session> session, const Bytes& body) {
            try {
                auto j = json::parse(std::string(body.begin(), body.end()));
                auto auth_request = j.get<AuthRequest>();
                
                auto token = auth_service->authenticate(auth_request.email, auth_request.password);
                
                json response = token;
                session->close(OK, response.dump(), {{"Content-Type", "application/json"}});
            } catch (const std::exception& e) {
                json error = {{"error", e.what()}};
                session->close(UNAUTHORIZED, error.dump(), {{"Content-Type", "application/json"}});
            }
        });
    }

    void verify_token_handler(const std::shared_ptr<Session> session) {
        const auto request = session->get_request();
        auto token = request->get_header("Authorization");
        
        if (token.empty()) {
            json error = {{"error", "No token provided"}};
            session->close(BAD_REQUEST, error.dump(), {{"Content-Type", "application/json"}});
            return;
        }

        // Remove "Bearer " prefix if present
        if (token.substr(0, 7) == "Bearer ") {
            token = token.substr(7);
        }

        bool is_valid = auth_service->verify_token(token);
        json response = {{"valid", is_valid}};
        session->close(OK, response.dump(), {{"Content-Type", "application/json"}});
    }

    void check_permission_handler(const std::shared_ptr<Session> session) {
        const auto request = session->get_request();
        size_t content_length = request->get_header("Content-Length", 0);

        session->fetch(content_length, [](const std::shared_ptr<Session> session, const Bytes& body) {
            try {
                auto j = json::parse(std::string(body.begin(), body.end()));
                auto check = j.get<PermissionCheck>();
                
                bool has_permission = auth_service->check_permission(
                    check.user_id,
                    check.resource,
                    check.action
                );
                
                json response = {{"has_permission", has_permission}};
                session->close(OK, response.dump(), {{"Content-Type", "application/json"}});
            } catch (const std::exception& e) {
                json error = {{"error", e.what()}};
                session->close(BAD_REQUEST, error.dump(), {{"Content-Type", "application/json"}});
            }
        });
    }
}

int main() {
    // Initialize services
    std::string db_url = std::getenv("DATABASE_URL");
    if (db_url.empty()) {
        db_url = "postgresql://testuser:testpass@postgres:5432/testingdb";
    }

    auto db = std::make_unique<Database>(db_url);
    auth_service = std::make_unique<AuthService>(std::move(db));

    // Setup API endpoints
    auto authenticate = std::make_shared<Resource>();
    authenticate->set_path("/auth/login");
    authenticate->set_method_handler("POST", authenticate_handler);

    auto verify = std::make_shared<Resource>();
    verify->set_path("/auth/verify");
    verify->set_method_handler("GET", verify_token_handler);

    auto check_permission = std::make_shared<Resource>();
    check_permission->set_path("/auth/check-permission");
    check_permission->set_method_handler("POST", check_permission_handler);

    // Setup and start service
    auto settings = std::make_shared<Settings>();
    settings->set_port(9000);
    settings->set_default_header("Connection", "close");

    Service service;
    service.publish(authenticate);
    service.publish(verify);
    service.publish(check_permission);
    service.start(settings);

    return EXIT_SUCCESS;
}
