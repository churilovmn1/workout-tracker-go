# Workout Tracker

Сервис для отслеживания тренировок. Пользователи логируют тренировки, отслеживают упражнения, видят личные рекорды и прогресс. Тренеры ведут расписание клиентов и оставляют комментарии к тренировкам. Доступ через веб-интерфейс, REST API и Telegram-бота.

## Стек

- **Бэкенд:** Go (chi router, pgx/v5)
- **База данных:** PostgreSQL
- **Брокер сообщений:** Redis (очередь уведомлений)
- **Авторизация:** JWT (HS256), пароли — bcrypt
- **Фронтенд:** HTML/JS (SPA)
- **Бот:** Telegram (long-polling)

## Архитектура

Запрос проходит сверху вниз по трём слоям. Каждый слой зависит только от
интерфейса нижележащего — зависимости собираются в `cmd/main.go`, без глобального
состояния и DI-фреймворка.

```
HTTP / Telegram
      │
      ▼
  handler/      — разбор запроса, JWT-аутентификация, коды ответов
      │           (AuthMiddleware → user_id+role в контекст; AdminOnly; rate limit)
      ▼
  service/      — бизнес-логика, проверки прав владельца, JWT, хеширование
      │           (каждый сервис объявляет свой узкий интерфейс репозитория)
      ▼
  repository/   — SQL через пул pgx; одна сущность на файл
      │
      ▼
  PostgreSQL
```

Асинхронные уведомления идут в обход HTTP: handler/service публикуют события в
Redis-очередь (`internal/broker`), а воркер в отдельной горутине читает очередь и
шлёт сообщения в Telegram (комментарий тренера, напоминание за час до тренировки).

```
cmd/main.go        — сборка зависимостей, запуск HTTP-сервера, бота и воркера
config/            — чтение переменных окружения
internal/
  models/          — структуры данных
  repository/      — слой доступа к БД (pgx/v5)
  service/         — бизнес-логика
  handler/         — chi-роутер, middleware, HTTP-хэндлеры
  broker/          — Redis-очередь: Publisher, Subscriber, Worker, события
bot/               — Telegram-бот
migrations/        — SQL-миграции (golang-migrate)
web/               — HTML-шаблон + статика (SPA)
```

## Сложность алгоритмов

Оценки для основных операций (`n` — число строк, по которым идёт выборка/обход):

| Операция                          | Сложность | Примечание |
|-----------------------------------|-----------|------------|
| `GetByID` (user/exercise/workout) | O(1)      | поиск по PRIMARY KEY / индексу |
| `GetByTelegramID`                 | O(1)      | по индексу `idx_users_telegram_id` |
| `ListByUser` (тренировки)         | O(n)      | по индексу `idx_workouts_user_id`, сортировка по date |
| `List` упражнений (фильтр)        | O(n)      | последовательное сканирование с фильтром по группе/тексту |
| `GetPersonalRecords`              | O(n)      | `DISTINCT ON` по exercise_id после индексной выборки |
| `GetWeeklyVolume`                 | O(n)      | агрегатная сумма по тренировкам за 7 дней |
| Rate limiter (на запрос)          | O(1)      | амортизированно, `sync.Map` по IP |
| Брокер: публикация/чтение события | O(1)      | Redis `LPUSH` / `BRPOP` |
| Брокер: скан напоминаний (раз/мин)| O(k)      | `k` — число сессий в окне «через час» |
| `calcStreak` (бот)                | O(n)      | один проход по списку тренировок |

## Как запустить

### Полный стек через Docker Compose

```bash
docker compose up -d              # поднимает db + redis + app
docker compose run --rm migrate   # применить миграции
```

### Локально

```bash
cp .env.example .env              # задать DATABASE_URL, JWT_SECRET, при желании REDIS_URL и TELEGRAM_BOT_TOKEN
docker compose up -d db redis     # только зависимости
make migrate-up                   # применить миграции
make run                          # go run ./cmd/
```

### Команды

```bash
make build        # сборка в bin/workout-tracker
make run          # go run ./cmd/
make test         # go test ./... -v
make lint         # golangci-lint run ./...
make migrate-up   # применить все миграции
make migrate-down # откатить последнюю миграцию

# Бенчмарки репозитория (нужен доступ к БД, иначе пропускаются):
DATABASE_URL=postgres://postgres:postgres@localhost:5433/workout_tracker?sslmode=disable \
  go test ./internal/repository/... -bench=. -run='^$'
```

## API

| Метод  | Путь                                | Описание                  | Роль  |
|--------|-------------------------------------|---------------------------|-------|
| POST   | /api/auth/register                  | Регистрация               | -     |
| POST   | /api/auth/login                     | Вход                      | -     |
| GET    | /api/exercises                      | Список упражнений         | user  |
| POST   | /api/exercises                      | Создать упражнение        | admin |
| PUT    | /api/exercises/:id                  | Обновить упражнение       | admin |
| DELETE | /api/exercises/:id                  | Удалить упражнение        | admin |
| GET    | /api/workouts                       | Список тренировок         | user  |
| POST   | /api/workouts                       | Создать тренировку        | user  |
| GET    | /api/workouts/:id                   | Детали тренировки         | user  |
| PUT    | /api/workouts/:id                   | Обновить тренировку       | user  |
| DELETE | /api/workouts/:id                   | Удалить тренировку        | user  |
| POST   | /api/workouts/:id/copy              | Скопировать тренировку    | user  |
| GET    | /api/stats/pr                       | Личные рекорды            | user  |
| GET    | /api/stats/volume                   | Объём за неделю           | user  |
| GET    | /api/templates                      | Список шаблонов           | user  |
| POST   | /api/templates                      | Создать шаблон            | user  |
| GET    | /api/admin/users                    | Список клиентов           | admin |
| GET    | /api/admin/users/:id/workouts       | Тренировки клиента        | admin |
| PUT    | /api/admin/workouts/:id/comment     | Комментарий тренера       | admin |
| GET    | /api/admin/schedule                 | Расписание на неделю      | admin |
| POST   | /api/admin/schedule                 | Добавить запись           | admin |
| GET    | /debug/pprof/                       | Профилировщик pprof       | admin |

## Переменные окружения

| Переменная         | Обязательна | Описание                                            | По умолчанию              |
|--------------------|-------------|-----------------------------------------------------|---------------------------|
| DATABASE_URL       | да          | URL подключения к PostgreSQL                        | -                         |
| JWT_SECRET         | нет         | Секрет для подписи JWT                              | `default-secret-change-me`|
| PORT               | нет         | Порт HTTP-сервера                                   | `8080`                    |
| TELEGRAM_BOT_TOKEN | нет         | Токен Telegram-бота (пусто → бот выключен)          | -                         |
| REDIS_URL          | нет         | URL Redis (пусто → брокер выключен, публикация no-op)| -                         |
