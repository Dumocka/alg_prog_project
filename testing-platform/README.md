# Testing Platform

Distributed application for conducting mass surveys and testing.

## Project Structure

- `core/` - Main testing logic module (Go)
- `auth/` - Authorization module (C++)
- `telegram-client/` - Telegram bot client (Python)
- `web-client/` - Web interface (JavaScript/React)

## Requirements

- Go 1.20+
- C++ 17
- Python 3.9+
- Node.js 18+
- PostgreSQL 14+
- Redis

## Architecture

The application is designed as a microservices architecture with the following components:

1. Core Module (Go):
   - Resource management (users, tests, questions, answers)
   - Access control
   - Business logic implementation
   - REST API

2. Auth Module (C++):
   - User rights management
   - Permission validation
   - Authentication proxy
   - OAuth2 integration

3. Telegram Client (Python):
   - Telegram Bot API integration
   - User interaction handling
   - Content adaptation for Telegram

4. Web Client (JavaScript/React):
   - Single Page Application
   - Responsive UI
   - Real-time updates

## Communication

Services communicate via REST APIs and message queues (RabbitMQ) for asynchronous operations.

## Deployment

The application supports horizontal scaling through containerization (Docker) and orchestration (Kubernetes).
