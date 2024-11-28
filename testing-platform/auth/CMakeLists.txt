cmake_minimum_required(VERSION 3.10)
project(auth_prototype)

set(CMAKE_CXX_STANDARD 17)

# TODO: добавить поиск зависимостей
# find_package(OpenSSL REQUIRED)
# find_package(CURL REQUIRED)

add_library(${PROJECT_NAME}
    src/auth_service.cpp
)

target_include_directories(${PROJECT_NAME}
    PUBLIC
        ${CMAKE_CURRENT_SOURCE_DIR}/include
)

# TODO: добавить линковку с библиотеками
# target_link_libraries(${PROJECT_NAME}
#     PRIVATE
#         OpenSSL::SSL
#         OpenSSL::Crypto
#         CURL::libcurl
# )
