#include "oauth_provider.hpp"
#include <curl/curl.h>
#include <nlohmann/json.hpp>
#include <sstream>

namespace auth {

namespace {
    size_t WriteCallback(void* contents, size_t size, size_t nmemb, std::string* userp) {
        userp->append((char*)contents, size * nmemb);
        return size * nmemb;
    }

    std::string make_http_request(const std::string& url, const std::string& method = "GET",
                                const std::string& data = "", const std::vector<std::string>& headers = {}) {
        CURL* curl = curl_easy_init();
        std::string response;

        if (curl) {
            curl_easy_setopt(curl, CURLOPT_URL, url.c_str());
            curl_easy_setopt(curl, CURLOPT_WRITEFUNCTION, WriteCallback);
            curl_easy_setopt(curl, CURLOPT_WRITEDATA, &response);

            if (method == "POST") {
                curl_easy_setopt(curl, CURLOPT_POST, 1L);
                if (!data.empty()) {
                    curl_easy_setopt(curl, CURLOPT_POSTFIELDS, data.c_str());
                }
            }

            struct curl_slist* header_list = nullptr;
            for (const auto& header : headers) {
                header_list = curl_slist_append(header_list, header.c_str());
            }
            if (header_list) {
                curl_easy_setopt(curl, CURLOPT_HTTPHEADER, header_list);
            }

            CURLcode res = curl_easy_perform(curl);
            if (res != CURLE_OK) {
                throw std::runtime_error(std::string("curl_easy_perform() failed: ") + 
                                      curl_easy_strerror(res));
            }

            if (header_list) {
                curl_slist_free_all(header_list);
            }
            curl_easy_cleanup(curl);
        }

        return response;
    }
}

// GitHub OAuth Implementation
GitHubOAuth::GitHubOAuth(const std::string& client_id, const std::string& client_secret)
    : client_id_(client_id), client_secret_(client_secret) {}

std::string GitHubOAuth::get_auth_url(const std::string& state) const {
    std::stringstream ss;
    ss << "https://github.com/login/oauth/authorize"
       << "?client_id=" << client_id_
       << "&state=" << state
       << "&scope=user:email";
    return ss.str();
}

User GitHubOAuth::exchange_code(const std::string& code) {
    // Exchange code for access token
    std::stringstream token_url;
    token_url << "https://github.com/login/oauth/access_token"
              << "?client_id=" << client_id_
              << "&client_secret=" << client_secret_
              << "&code=" << code;

    std::vector<std::string> headers = {"Accept: application/json"};
    std::string token_response = make_http_request(token_url.str(), "POST", "", headers);
    auto token_json = nlohmann::json::parse(token_response);
    
    if (token_json.contains("error")) {
        throw std::runtime_error("GitHub OAuth error: " + token_json["error"].get<std::string>());
    }

    std::string access_token = token_json["access_token"];

    // Get user info
    headers = {
        "Accept: application/json",
        "Authorization: token " + access_token
    };
    std::string user_response = make_http_request("https://api.github.com/user", "GET", "", headers);
    auto user_json = nlohmann::json::parse(user_response);

    User user;
    user.email = user_json["email"];
    user.name = user_json["name"];
    // Note: In a real implementation, you would need to generate or fetch a user ID
    user.id = 0;  // This should be set by the database
    user.created_at = ""; // This should be set by the database

    return user;
}

// Yandex OAuth Implementation
YandexOAuth::YandexOAuth(const std::string& client_id, const std::string& client_secret)
    : client_id_(client_id), client_secret_(client_secret) {}

std::string YandexOAuth::get_auth_url(const std::string& state) const {
    std::stringstream ss;
    ss << "https://oauth.yandex.ru/authorize"
       << "?response_type=code"
       << "&client_id=" << client_id_
       << "&state=" << state;
    return ss.str();
}

User YandexOAuth::exchange_code(const std::string& code) {
    // Exchange code for access token
    std::stringstream token_url;
    token_url << "https://oauth.yandex.ru/token"
              << "?grant_type=authorization_code"
              << "&code=" << code
              << "&client_id=" << client_id_
              << "&client_secret=" << client_secret_;

    std::string token_response = make_http_request(token_url.str(), "POST");
    auto token_json = nlohmann::json::parse(token_response);
    
    if (token_json.contains("error")) {
        throw std::runtime_error("Yandex OAuth error: " + token_json["error"].get<std::string>());
    }

    std::string access_token = token_json["access_token"];

    // Get user info
    std::vector<std::string> headers = {
        "Authorization: OAuth " + access_token
    };
    std::string user_response = make_http_request("https://login.yandex.ru/info", "GET", "", headers);
    auto user_json = nlohmann::json::parse(user_response);

    User user;
    user.email = user_json["default_email"];
    user.name = user_json["real_name"];
    // Note: In a real implementation, you would need to generate or fetch a user ID
    user.id = 0;  // This should be set by the database
    user.created_at = ""; // This should be set by the database

    return user;
}

} // namespace auth
