cmake_minimum_required(VERSION 3.15)
project(auth VERSION 1.0.0 LANGUAGES CXX)

set(CMAKE_CXX_STANDARD 17)
set(CMAKE_CXX_STANDARD_REQUIRED ON)

# Опции сборки
option(BUILD_TESTS "Build tests" ON)
option(BUILD_EXAMPLES "Build examples" OFF)

# Зависимости
find_package(OpenSSL REQUIRED)
find_package(CURL REQUIRED)
find_package(nlohmann_json REQUIRED)
find_package(cpp_redis REQUIRED)
find_package(jwt-cpp REQUIRED)
find_package(libpqxx REQUIRED)

# Основная библиотека
add_library(auth
    src/authorization_server.cpp
    src/jwt_manager.cpp
    src/code_authentication.cpp
    src/oauth_provider.cpp
    src/database.cpp
)

target_include_directories(auth
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
    PRIVATE
        ${CMAKE_CURRENT_SOURCE_DIR}/src
)

target_link_libraries(auth
    PUBLIC
        OpenSSL::SSL
        OpenSSL::Crypto
        CURL::libcurl
        nlohmann_json::nlohmann_json
        cpp_redis
        jwt-cpp
        libpqxx::pqxx
)

# Тесты
if(BUILD_TESTS)
    enable_testing()
    add_subdirectory(tests)
endif()

# Примеры
if(BUILD_EXAMPLES)
    add_subdirectory(examples)
endif()

# Установка
install(TARGETS auth
    LIBRARY DESTINATION lib
    ARCHIVE DESTINATION lib
    RUNTIME DESTINATION bin
)

install(DIRECTORY include/
    DESTINATION include
    FILES_MATCHING PATTERN "*.hpp"
)
