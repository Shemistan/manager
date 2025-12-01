# Manager Service

Эталонный HTTP-сервис на Go с правильной слоистой архитектурой и чистым кодом.

## Архитектура

Сервис построен с соблюдением классической трёхслойной архитектуры:

### Структура проекта

```
manager/
├── cmd/
│   ├── manager/           # Точка входа сервиса
│   │   └── main.go
│   └── migrator/          # Точка входа мигратора
│       └── main.go
├── internal/
│   ├── api/               # HTTP-слой (обработчики и маршруты)
│   │   └── manager/
│   │       ├── handler.go
│   │       ├── handler_test.go
│   │       └── response.go
│   ├── app/               # Инициализация приложения
│   │   ├── manager/
│   │   │   └── app.go
│   │   └── migrator/
│   │       └── migrator.go
│   ├── config/            # Конфигурация
│   │   └── config.go
│   ├── service/           # Интерфейсы бизнес-логики
│   │   ├── service.go
│   │   └── manager/       # Реализация для manager
│   │       ├── service.go
│   │       └── service_test.go
│   └── storage/           # Интерфейсы работы с БД
│       ├── storage.go
│       └── manager/       # Реализация для PostgreSQL
│           └── storage.go
├── migration/             # SQL-миграции
│   └── 001_init_health.sql
├── app.toml               # Конфиг приложения
├── .env.example           # Пример переменных окружения
├── Dockerfile             # Docker-образ
├── docker-compose.yml     # Оркестрация контейнеров
├── go.mod
├── go.sum
└── README.md
```

## Принципы архитектуры

### 1. **Связь слоёв через интерфейсы**

- **API слой** → **Service слой**: через интерфейс `HealthService`
- **Service слой** → **Storage слой**: через интерфейс `HealthStorage`
- Нет прямых зависимостей между слоями
- Каждый слой зависит от абстракций, а не от конкретных реализаций

### 2. **Разделение ответственности**

- **API** (`internal/api/manager/`): обработка HTTP-запросов, валидация, сериализация ответов
- **Service** (`internal/service/manager/`): бизнес-логика, координация между слоями
- **Storage** (`internal/storage/manager/`): работа с БД, выполнение SQL-запросов
- **App** (`internal/app/manager/`): "склейка" зависимостей, инициализация

### 3. **Конфигурация**

- `app.toml`: общие настройки (не содержит секретов)
- `.env`: секреты локально (не коммитится)
- Переменные окружения для Docker

## Функциональность

### Endpoint: GET /health

```http
GET /health HTTP/1.1
```

**Ответ (200 OK):**
```json
{
  "status": "success"
}
```

**Поведение:**
- Возвращает JSON `{"status":"success"}` с HTTP 200
- При каждом вызове записывает информацию о вызове в таблицу `health_calls` PostgreSQL
- При ошибке БД возвращает HTTP 500 с JSON `{"status":"error","error":"..."}`

## Запуск локально (без Docker)

### Требования

- Go 1.23.4+
- PostgreSQL 16+ (запущен отдельно)

### Шаги

1. **Установить зависимости:**
   ```bash
   go mod download
   ```

2. **Создать БД и пользователя:**
   ```bash
   psql -U postgres -c "CREATE USER manager_user WITH PASSWORD 'secure_password_here';"
   psql -U postgres -c "CREATE DATABASE manager_db OWNER manager_user;"
   psql -U postgres -c "ALTER USER manager_user CREATEDB;"
   ```

3. **Создать `.env` файл:**
   ```bash
   cp .env.example .env
   ```

4. **Запустить миграции:**
   ```bash
   go run ./cmd/migrator/main.go
   ```
   Создаст таблицу `health_calls` в БД.

5. **Запустить сервис:**
   ```bash
   go run ./cmd/manager/main.go
   ```
   Сервис будет доступен на `http://localhost:8081`

6. **Проверить здоровье:**
   ```bash
   curl http://localhost:8081/health
   ```

### Запуск тестов

```bash
go test ./... -v
```

Тесты покрывают:
- `HealthService.HandleHealth()` с успешным и ошибочным сценариями
- HTTP-обработчик `/health` с различными состояниями
- Проверка Content-Type, HTTP-кодов и структуры ответов

## Запуск с Docker Compose

### Требования

- Docker
- Docker Compose

### Шаги

1. **Создать `.env` файл (если не существует):**
   ```bash
   cp .env.example .env
   ```

2. **Поднять сервис:**
   ```bash
   docker-compose up --build
   ```

   Это запустит:
   - PostgreSQL контейнер (`db`)
   - Мигратор (`migrator`) — применит миграции и завершит работу
   - Сервис `manager` — HTTP-сервер, зависящий от успешного завершения мигратора

3. **Проверить здоровье:**
   ```bash
   curl http://localhost:8081/health
   ```

4. **Остановить:**
   ```bash
   docker-compose down
   ```

5. **Удалить данные БД:**
   ```bash
   docker-compose down -v
   ```

## API Примеры

### Успешный запрос

```bash
curl -X GET http://localhost:8081/health
```

```json
{
  "status": "success"
}
```

### Проверка БД

```bash
psql -U manager_user -d manager_db -h localhost
```

```sql
SELECT COUNT(*), MAX(called_at) FROM health_calls;
```

## Миграции

Миграции находятся в папке `migration/` и выполняются в алфавитном порядке.

### `001_init_health.sql`

Создаёт таблицу для логирования вызовов health endpoint:

```sql
CREATE TABLE IF NOT EXISTS health_calls (
    id SERIAL PRIMARY KEY,
    called_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX IF NOT EXISTS idx_health_calls_called_at ON health_calls(called_at);
```

## Конфигурация

### `app.toml`

```toml
service_name = "manager"
service_env = "local"  # local / prod
http_port = 8081

[database]
host = "localhost"
port = 5432
user = "manager_user"
name = "manager_db"
sslmode = "disable"

[tls]
enabled = false
cert_file = ""
key_file = ""
ca_file = ""
```

### Переменные окружения

#### Database
```
DB_HOST=localhost              # PostgreSQL хост
DB_PORT=5432                   # PostgreSQL порт
DB_USER=manager_user           # Пользователь БД
DB_PASSWORD=secure_password_here # Пароль БД
DB_NAME=manager_db             # Имя БД
DB_SSLMODE=disable             # SSL режим (disable/require)
```

#### Application
```
APP_PORT=8081                  # Порт HTTP сервера
```

#### TLS (по умолчанию отключен)
```
TLS_ENABLED=false              # Включить TLS (true/false)
TLS_CERT_FILE=                 # Путь к сертификату сервера
TLS_KEY_FILE=                  # Путь к приватному ключу сервера
TLS_CA_FILE=                   # Путь к сертификату CA для проверки client cert
```

### Пример `.env` для локальной разработки

```bash
DB_USER=manager_user
DB_PASSWORD=secure_password_here
DB_SSLMODE=disable
```

## Зависимости

- `github.com/lib/pq` — PostgreSQL драйвер для Go
- `github.com/BurntSushi/toml` — парсер TOML конфигов
- `github.com/joho/godotenv` — загрузка переменных из `.env`

## Качество кода

- ✅ `go build ./...` проходит без ошибок
- ✅ `go test ./...` все тесты проходят
- ✅ Явная обработка ошибок
- ✅ Идиоматичный Go код
- ✅ Интерфейсы для зависимостей
- ✅ Unit тесты с мок-объектами

## TLS конфигурация (отключено по умолчанию)

Сервис готов к использованию mTLS (mutual TLS) аутентификации. В данный момент TLS отключен, но архитектура полностью поддерживает включение.

### Включение TLS

Для включения HTTPS и проверки сертификатов клиента:

1. **Подготовить сертификаты:**
   ```bash
   # Сертификат и ключ сервера
   openssl req -x509 -newkey rsa:4096 -keyout server.key -out server.crt -days 365 -nodes

   # CA сертификат для проверки client cert (опционально для mTLS)
   ```

2. **Установить переменные окружения:**
   ```bash
   export TLS_ENABLED=true
   export TLS_CERT_FILE=/path/to/server.crt
   export TLS_KEY_FILE=/path/to/server.key
   export TLS_CA_FILE=/path/to/ca.crt  # для проверки client cert
   ```

3. **Запустить сервис:**
   ```bash
   go run ./cmd/manager/main.go
   ```

### Параметры TLS

- **TLS_ENABLED**: Включить/выключить TLS (true/false, по умолчанию false)
- **TLS_CERT_FILE**: Путь к файлу с сертификатом сервера (*.crt)
- **TLS_KEY_FILE**: Путь к файлу с приватным ключом сервера (*.key)
- **TLS_CA_FILE**: Путь к файлу с CA сертификатом для проверки client cert (опционально)

## Развёртывание в продакшене

### Изменения для продакшена

1. **`app.toml`**: установить `service_env = "prod"`
2. **Миграции**: выполнить один раз перед стартом сервиса
3. **SSL БД**: изменить `sslmode` в конфиге на `"require"`
4. **Переменные окружения**: использовать безопасные хранилища (AWS Secrets Manager, Vault, и т.д.)
5. **Dockerfile**: убедиться в многостейджовой сборке и минимальном образе

## Проблемы и решения

### Ошибка: "failed to ping database"

Убедиться, что PostgreSQL запущен и доступен на `localhost:5432` с правильными креденшалами.

### Ошибка: "no such file or directory: app.toml"

Убедиться, что вы запускаете сервис из корневой директории проекта.

### Ошибка: "migration directory: no such file or directory"

Убедиться, что папка `migration/` существует и содержит файлы `.sql`.

## Лицензия

MIT
