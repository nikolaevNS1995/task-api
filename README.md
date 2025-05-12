
# task-api

## Описание проекта

`task-api` - это веб-приложение для управления задачами с возможностью регистрации, просмотра, обновления и удаления задач. Проект включает в себя REST API с JWT-аутентификацией и использует PostgreSQL для хранения данных.

## Основные возможности

- Аутентификация пользователей с использованием JWT
- Создание, редактирование и удаление задач
- Просмотр списка задач
- Простой интерфейс для взаимодействия с задачами через API

## Архитектура

Проект построен с использованием принципов чистой архитектуры и разделен на следующие слои:

- **Domain** - бизнес-логика и сущности
- **Infrastructure** - реализация хранилищ данных, API и сервисов
- **Application** - координация между слоями и конфигурация приложения

Для управления зависимостями используется библиотека `fx`.

## Технологический стек

- **Язык**: Go
- **Web-фреймворк**: Gin
- **База данных**: PostgreSQL
- **Миграции**: golang-migrate
- **Логирование**: Zap
- **Телеметрия**: OpenTelemetry
- **Управление зависимостями**: Uber FX
- **Аутентификация**: JWT

## Требования

- Go 1.24+
- PostgreSQL 14+

## Установка и запуск

### Подготовка окружения

1. Клонировать репозиторий:
```bash
git clone https://github.com/nikolaevNS1995/task-api.git
cd task-api
```

2. Создать и настроить файл конфигурации `.env`

3. Создать базу данных PostgreSQL:
```bash
psql -U postgres
CREATE DATABASE task_db;
```

### Запуск приложения

1. Установить зависимости:
```bash
go mod download
```

2. Запустить приложение:
```bash
go run cmd/app/main.go
```

Или с указанием пути к конфигурационному файлу:
```bash
PATH_CONFIG=./config/.env go run cmd/app/main.go
```

## Конфигурация

Приложение можно настроить через переменные окружения. Основные параметры конфигурации:

### Основные настройки
```
APP_NAME="task-app"                   # Имя приложения
ADDRESS_SERVER="0.0.0.0:8080"         # Адрес и порт сервера
APP_ENV="development"                 # Окружение (development, production)
```

### База данных
```
POSTGRES_DB_HOST="postgres"           # Хост базы данных
POSTGRES_DB_PORT=5432                 # Порт базы данных
POSTGRES_DB_USER="user"               # Пользователь базы данных
POSTGRES_DB_PASSWORD="password"       # Пароль базы данных
POSTGRES_DB_NAME="task_db"            # Имя базы данных
POSTGRES_DB_SSLMODE="disable"         # Режим SSL для подключения
```

### Аутентификация
```
JWT_SECRET=SecretKey                 # Секретный ключ для JWT
JWT_ALGORITHM="HS256"                # Алгоритм JWT
JWT_EXPIRY="60m"                     # Время жизни токена
JWT_REFRESH_EXPIRY="43200m"          # Время жизни refresh токена
```

### Логирование
```
LOGGER_LEVEL="debug"                 # Уровень логирования
LOGGER_OUTPUT="stdout"               # Вывод логов (stdout, stderr, файл)
LOGGER_FILE_PATH="./logs/app.log"    # Путь к файлу логов
LOGGER_MAX_SIZE_MB=10                # Максимальный размер файла лога
LOGGER_MAX_BACKUPS=5                 # Количество резервных копий логов
LOGGER_MAX_AGE_DAYS=30               # Время хранения логов
LOGGER_COMPRESS=true                 # Сжатие логов
LOGGER_APP_ENV="${APP_ENV}"          # Окружение для логирования
```

### Телеметрия
```
TELEMETRY_HOST="localhost"            # Адрес хоста телеметрии
TELEMETRY_PORT=4317                   # Порт телеметрии
TELEMETRY_LOCAL=true                  # Локальный режим телеметрии
```

## Переменные окружения

Приложение поддерживает следующие переменные окружения:

### Основные настройки
```
PATH_CONFIG             # Путь к файлу конфигурации
ADDRESS_SERVER          # Адрес и порт сервиса (например, ":8080")
```

### Логирование
```
LOGGER_LEVEL            # Уровень логирования (debug, info, warn, error)
LOGGER_OUTPUT           # Вывод логов (stdout, stderr, файл)
LOGGER_FILE_PATH       # Путь к файлу логов
LOGGER_MAX_SIZE_MB     # Максимальный размер файла лога (в МБ)
LOGGER_MAX_BACKUPS     # Количество резервных копий логов
LOGGER_MAX_AGE_DAYS    # Время хранения логов (в днях)
LOGGER_COMPRESS        # Сжатие логов (true/false)
```

### Телеметрия
```
TELEMETRY_HOST          # Адрес хоста телеметрии
TELEMETRY_PORT          # Порт хоста телеметрии
TELEMETRY_LOCAL         # Локальный режим телеметрии (true/false)
```

### База данных
```
POSTGRES_DB_HOST       # Хост базы данных
POSTGRES_DB_PORT       # Порт базы данных
POSTGRES_DB_USER       # Пользователь базы данных
POSTGRES_DB_PASSWORD   # Пароль базы данных
POSTGRES_DB_NAME       # Имя базы данных
POSTGRES_DB_SSLMODE    # Режим SSL для подключения к базе данных
```

## Структура проекта

```
task-api/
├── cmd/                    # Точка входа в приложение
├── internal/               # Внутренний код приложения
│   ├── adapters/           # Адаптеры для внешних сервисов
│   ├── app/                # Конфигурация приложения
│   ├── constants/          # Константы
│   ├── domain/             # Бизнес-логика
│   │   ├── entities/       # Сущности
│   │   └── usecases/       # Сценарии использования
│   └── infrastructure/     # Инфраструктурный код
│       ├── api/            # API интерфейсы
│       ├── repository/     # Репозитории для работы с данными
│       └── service/        # Сервисы
├── migrations/             # Миграции базы данных
└── pkg/                    # Общие пакеты
    ├── config/             # Конфигурация
    ├── connectors/         # Коннекторы к внешним сервисам
    ├── telemetry/          # Телеметрия
    └── utils/              # Вспомогательные функции
```

## API Endpoints

### Задачи
- `GET /v1/tasks` - Получение списка задач
- `POST /v1/tasks` - Создание новой задачи
- `PUT /v1/tasks/{id}` - Обновление задачи
- `DELETE /v1/tasks/{id}` - Удаление задачи