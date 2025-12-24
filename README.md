# Subscription Service API 🐳

Микросервис для управления подписками с полной Docker-инфраструктурой.

## 🚀 Быстрый запуск через Docker

```bash

# 1. Клонируй репозиторий
git clone https://github.com/Negst1/subscribe_project.git

# 2. Запустите всё одной командой
docker-compose up --build -d

# 3. Пример .env

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=subscriptions_db
SERVER_PORT=8080

#4. Данные от pgAdmin

Логин: admin@sub.com
Пароль: admin

Connection: 
Host: postgres
Port: 5432
Username: postgres
Pass: postgres