# Подготовка к защите проекта Workout Tracker

---

## 1. Архитектура — зачем нужен каждый слой

Проект разделён на три слоя. Каждый знает только о следующем вниз — это называется **разделение ответственности**.

```
Запрос от браузера
       │
  [ Handler ]   ← принимает HTTP-запрос, достаёт параметры, отвечает JSON
       │
  [ Service ]   ← бизнес-логика: «только владелец может удалить свою тренировку»
       │
 [ Repository ] ← SQL-запросы к базе данных
       │
  [ PostgreSQL ]
```

### Handler (`internal/handler/`)
**Что делает:** получает HTTP-запрос, читает JSON из тела, достаёт `user_id` из context, вызывает нужный сервис, возвращает JSON-ответ с правильным HTTP-кодом.

**Чем НЕ занимается:** не знает о SQL, не проверяет права владельца — это дело сервиса.

Пример: `handler/workout.go` — метод `Create` читает JSON, формирует структуру `Workout`, передаёт в `workoutService.Create`, возвращает `201 Created`.

### Service (`internal/service/`)
**Что делает:** бизнес-правила. Например, `WorkoutService.GetByID` сначала берёт тренировку из БД, затем **проверяет** `w.UserID != userID` — если чужая, возвращает `ErrForbidden`.

**Почему сервис не знает про HTTP:** если завтра добавить CLI или gRPC, сервис останется тем же.

### Repository (`internal/repository/`)
**Что делает:** только SQL. Один файл — одна таблица. Никакой бизнес-логики.

**Пример:** `WorkoutRepository.Create` выполняет транзакцию: вставляет `workouts`, потом все `workout_exercises`. Если любой INSERT упал — `tx.Rollback`. Это атомарность.

---

## 2. CRUD — где реализован

CRUD = Create / Read / Update / Delete.

### Упражнения (Exercises) — полный CRUD

| Операция | HTTP | Файл | Метод сервиса |
|----------|------|------|---------------|
| Create | `POST /api/exercises` | `handler/exercise.go` | `ExerciseService.Create` |
| Read (список) | `GET /api/exercises` | `handler/exercise.go` | `ExerciseService.List` |
| Read (один) | `GET /api/exercises/{id}` | `handler/exercise.go` | `ExerciseService.GetByID` |
| Update | `PUT /api/exercises/{id}` | `handler/exercise.go` | `ExerciseService.Update` |
| Delete | `DELETE /api/exercises/{id}` | `handler/exercise.go` | `ExerciseService.Delete` |

### Тренировки (Workouts) — полный CRUD + копирование

| Операция | HTTP | Файл |
|----------|------|------|
| Create | `POST /api/workouts` | `handler/workout.go` → `WorkoutService.Create` |
| Read (список) | `GET /api/workouts` | `handler/workout.go` → `WorkoutService.ListByUser` |
| Read (один) | `GET /api/workouts/{id}` | `handler/workout.go` → `WorkoutService.GetByID` |
| Update | `PUT /api/workouts/{id}` | `handler/workout.go` → `WorkoutService.Update` |
| Delete | `DELETE /api/workouts/{id}` | `handler/workout.go` → `WorkoutService.Delete` |
| Copy | `POST /api/workouts/{id}/copy` | `handler/workout.go` → `WorkoutService.CopyWorkout` |

### Остальные сущности

| Сущность | CRUD в | Файл репозитория |
|----------|--------|-----------------|
| Пользователи | `handler/auth.go`, `handler/admin.go` | `repository/user.go` |
| Шаблоны | `handler/template.go` | `repository/template.go` |
| Расписание | `handler/admin.go` | `repository/schedule.go` |
| Метрики тела | `handler/metrics.go` | `repository/metrics.go` |

---

## 3. Темы курса — где в проекте

---

### SOLID

Все пять принципов применены:

**S — Single Responsibility (одна ответственность)**
Каждый файл делает одно дело: `UserRepository` — только SQL по таблице users, `AuthService` — только аутентификация, `WorkoutHandler` — только HTTP для тренировок.

**O — Open/Closed (открыт для расширения, закрыт для изменения)**
Сервисы работают через интерфейсы. Чтобы добавить кэширование — создай новый тип, реализующий тот же интерфейс, без изменения сервиса.

**L — Liskov Substitution (подстановка Лисков)**
`repository.WorkoutRepository` удовлетворяет интерфейсу `workoutRepository` из service-пакета. В тестах вместо него передаётся `mockWorkoutRepo` — и сервис не замечает разницы (`internal/service/workout_test.go`).

**I — Interface Segregation (разделение интерфейсов)**
Каждый сервис объявляет свой **узкий** интерфейс с только нужными методами. `AdminService` объявляет `adminUserRepository` с одним методом `List` — ему не нужен `GetByLogin`. Файл `internal/service/admin.go`, строки 10–26.

**D — Dependency Inversion (инверсия зависимостей)**
Сервисы зависят от интерфейсов, а не от конкретных структур. Конкретные репозитории передаются сверху через конструктор в `cmd/main.go`. Это и есть Dependency Injection без фреймворка.

---

### Дженерики (Generics)

В Go 1.18+ появились дженерики. В проекте используются через библиотеку pgx/v5.

**Где:** `internal/repository/user.go`, метод `List` — функция `pgx.CollectRows`:

```go
pgx.CollectRows(rows, func(row pgx.CollectableRow) (models.User, error) {
    var u models.User
    err := row.Scan(&u.ID, &u.Login, ...)
    return u, err
})
```

`pgx.CollectRows[T any]` — обобщённая функция, работающая с любым типом T. В разных репозиториях передаём `models.User`, `models.Workout`, `models.Exercise` — один код, разные типы.

Без дженериков пришлось бы писать отдельную функцию для каждого типа или возвращать `[]interface{}`.

---

### Указатели

**Где в проекте:**

- Все репозитории и сервисы передаются как указатели: `*WorkoutRepository`, `*AuthService` — чтобы не копировать структуру с пулом соединений при каждом вызове.
- `*models.User`, `*models.Workout` — возвращаем указатели из методов, `nil` означает «не найдено».
- `*pgxpool.Pool` — единственный пул, передаётся во все репозитории по указателю.

**Почему не значение:**
Если передать `pgxpool.Pool` по значению — скопируется вся структура со всеми соединениями. Каждый репозиторий будет работать со своей копией. Это не то, что нужно — нам нужен ОДИН общий пул.

**Пример:** `repository/user.go`:
```go
type UserRepository struct {
    pool *pgxpool.Pool  // указатель — все репозитории смотрят на один пул
}
```

---

### Graceful Shutdown (Beaty Close)

**Что это:** сервер не убивается мгновенно при нажатии Ctrl+C, а сначала дожидается завершения текущих запросов.

**Где:** `cmd/main.go`, строки в конце:

```go
quit := make(chan os.Signal, 1)
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
<-quit  // ждём сигнала

shutdownCtx, _ := context.WithTimeout(context.Background(), 5*time.Second)
srv.Shutdown(shutdownCtx)  // даём 5 секунд на завершение запросов
```

**Почему важно:** если сервер убить без graceful shutdown, клиент получит обрыв соединения вместо нормального ответа. В базе может остаться незакоммиченная транзакция.

---

### Брокер сообщений

**В нашем проекте** брокер был реализован на Redis (паттерн LPUSH/BRPOP) и удалён из итоговой версии. Но концепция важна для ответа на вопросы:

**Что это:** паттерн Publisher-Subscriber. Один компонент публикует событие (например, «тренировка создана»), другой асинхронно его обрабатывает (отправляет уведомление). Они не знают друг о друге.

**Зачем:** если уведомление падает, оно не роняет основной запрос. Обработка идёт в фоне.

**Паттерн в коде:** интерфейс `Publisher` с методом `Publish` — сервер не знал, реальный это Redis или заглушка `NoopPublisher`.

---

### make

В Go `make` — встроенная функция для инициализации слайсов, map и каналов.

**В проекте:**

```go
// service/workout.go — создаём слайс заданного размера
dst.Exercises = make([]models.WorkoutExercise, len(source.Exercises))

// handler/ratelimit.go — sync.Map инициализируется автоматически

// cmd/main.go — канал для сигналов ОС
quit := make(chan os.Signal, 1)
```

**Отличие от `new`:** `make` возвращает инициализированный тип (готовый к работе), `new` — указатель на нулевое значение. Для map, slice, chan всегда нужен `make`.

**Также:** `Makefile` в корне проекта — утилита для запуска команд сборки:
```
make build   → go build -o bin/workout-tracker ./cmd/
make test    → go test ./... -v
```

---

### Структура проекта

Стандартная Go-раскладка по соглашению сообщества:

```
cmd/          — точка входа (main.go). Только сборка зависимостей.
config/       — чтение окружения. Один файл, одна задача.
internal/     — приватный код: недоступен внешним модулям
  models/     — структуры данных (User, Workout, Exercise...)
  repository/ — SQL-запросы, один файл = одна таблица
  service/    — бизнес-логика
  handler/    — HTTP-обработчики, middleware
migrations/   — SQL-миграции, пронумерованы 001–014
web/          — фронтенд: HTML-шаблон + JS + CSS
docs/         — Swagger-документация (генерируется автоматически)
```

**Ключевое правило:** `internal/` защищён компилятором — внешний модуль не может импортировать этот код. Это намеренная инкапсуляция.

---

### HTTPS

**В проекте:** сервер запускает HTTP (`srv.ListenAndServe`). HTTPS в production обеспечивает reverse proxy (nginx, Caddy) перед приложением — это стандартная практика.

**Защита данных в проекте:**
- Пароли никогда не передаются в открытом виде после регистрации — только JWT-токен
- JWT подписан HMAC-SHA256 — подделать без секрета невозможно
- SQL-запросы параметризованы — исключена SQL-инъекция

---

### WebAssembly

В нашем проекте не используется. Go умеет компилироваться в WebAssembly (`GOOS=js GOARCH=wasm`), что позволяет запускать Go-код прямо в браузере. В нашем случае фронтенд — обычный JS, бэкенд — Go-сервер.

---

### GUI

**В проекте:** SPA (Single Page Application) на Vanilla JS.

- `web/templates/index.html` — HTML-шаблон, отдаётся сервером
- `web/static/app.js` — весь фронтенд-код: AJAX-запросы к API, отрисовка данных
- `web/static/style.css` — стили
- Графики — Chart.js (CDN)

Сервер отдаёт статику: `handler/web.go`, роут `r.Handle("/static/*", ...)`.

---

### Хранение паролей

**Принцип:** пароль никогда не хранится в открытом виде. Хранится только bcrypt-хэш.

**Где:** `internal/service/auth.go`:

```go
// При регистрации:
hash, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
// В БД записывается строка вида: $2a$10$...xyz...

// При входе:
bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
// Сравниваем пароль с хэшем, не расшифровывая
```

**Почему bcrypt, а не MD5/SHA256:**
1. Намеренно медленный — брутфорс занимает годы, а не часы
2. Встроена «соль» (salt) — два одинаковых пароля дают разные хэши
3. `DefaultCost = 10` — можно поднять сложность без изменения кода

---

### Return как можно раньше (Early Return / Guard Clauses)

**Принцип:** сначала проверяем ошибки, возвращаемся при первой проблеме. Не вкладываем успешный путь в `else`.

**Где в проекте** — весь `handler/workout.go`:

```go
// ПО ПРИНЦИПУ EARLY RETURN:
id, err := strconv.Atoi(chi.URLParam(r, "id"))
if err != nil {
    writeError(w, http.StatusBadRequest, "invalid workout id")
    return  // ← выходим сразу
}

workout, err := h.workoutService.GetByID(...)
if err != nil {
    if errors.Is(err, service.ErrForbidden) {
        writeError(w, http.StatusForbidden, "access denied")
        return  // ← выходим сразу
    }
    writeError(w, http.StatusNotFound, "workout not found")
    return  // ← выходим сразу
}

writeJSON(w, http.StatusOK, workout)  // ← счастливый путь в конце, без else
```

**Почему это лучше:** код читается сверху вниз. Видно все ошибочные случаи до основной логики. Не нужно искать глазами закрывающие скобки вложенных `if-else`.

---

### GIT (углуб)

**В проекте:**
```
git log --oneline:
c763e3b feat: body metrics, progress charts, exercise names in PR
ec3045f feat: redis broker, rate limiting, pprof, benchmarks, seed data
db2be94 feat: workout tracker - Go, PostgreSQL, JWT, Telegram bot, trainer panel
```

**Что важно знать:**
- Коммиты атомарны — каждый фиксирует одно логическое изменение
- `.gitignore` исключает `.env` (содержит секреты) и `bin/` (бинарник не нужен в репозитории)
- `main` — единственная ветка, для учебного проекта достаточно

---

### Сложность вычислений O(n)

| Операция | Сложность | Почему |
|----------|-----------|--------|
| `GetByID` (user, workout, exercise) | **O(1)** | Поиск по PRIMARY KEY — B-Tree индекс |
| `GetByLogin` | **O(1)** | Уникальный индекс по полю login |
| `ListByUser` (тренировки) | **O(n)** | Сканирование всех тренировок пользователя по `idx_workouts_user_id` |
| `GetPersonalRecords` | **O(n)** | `DISTINCT ON (exercise_id)` — один проход по всем упражнениям пользователя |
| `GetWeeklyVolume` | **O(n)** | `SUM(sets * reps * weight)` по тренировкам за 7 дней |
| Rate Limiter (на запрос) | **O(1)** | Амортизированно: `sync.Map` по IP-адресу |
| `WithSearch` фильтр упражнений | **O(n)** | `ILIKE` — последовательный поиск без full-text индекса |

**Файл:** `internal/repository/workout.go` — методы `GetPersonalRecords`, `GetWeeklyVolume`, `GetExerciseProgress`.

---

### Runtime

Go Runtime — среда выполнения, встроенная в каждый скомпилированный бинарник.

**Что даёт Runtime:**

1. **Сборщик мусора (GC)** — автоматически освобождает неиспользуемую память. Мы не вызываем `free()`.
2. **Горутины (goroutines)** — легковесные потоки. В `cmd/main.go`:
```go
go func() {
    srv.ListenAndServe()  // HTTP-сервер в горутине
}()
// главная горутина ждёт сигнала ОС
```
3. **Scheduler** — Runtime сам распределяет горутины по ядрам CPU (GOMAXPROCS).

**Горутины в проекте:** HTTP-сервер запускается в отдельной горутине, чтобы главная могла ждать сигнала завершения.

---

### godoc

**Что это:** стандарт документирования кода в Go. Комментарий перед функцией — это её документация.

**Правило:** комментарий начинается с имени функции.

```go
// NewPostgresPool creates a new connection pool to PostgreSQL.
func NewPostgresPool(ctx context.Context, databaseURL string) (*pgxpool.Pool, error) {
```

**В проекте:**
- Все публичные типы и функции задокументированы
- Swagger-аннотации (`// @Summary`, `// @Tags`) — расширение godoc для HTTP API
- Запустить документацию: `go doc ./internal/service/...`

---

### Профилировка (Profiling)

**Что это:** инструмент для поиска узких мест — что тормозит, где утекает память.

**Где в проекте:** `internal/handler/router.go` — маршруты `/debug/pprof/`:

```go
r.Route("/debug/pprof", func(r chi.Router) {
    r.Use(AuthMiddleware(authService))
    r.Use(AdminOnly)  // только для admin!
    r.HandleFunc("/", pprof.Index)
    r.HandleFunc("/profile", pprof.Profile)   // CPU профиль
    r.HandleFunc("/{profile}", pprof.Index)   // heap, goroutine, allocs...
})
```

**Как использовать:**
```bash
# CPU профиль на 30 секунд:
go tool pprof http://localhost:8080/debug/pprof/profile?seconds=30

# Анализ памяти:
go tool pprof http://localhost:8080/debug/pprof/heap
```

**Почему защищён AdminOnly:** профилировщик раскрывает внутренности приложения.

---

### Тестирование

**Где:** `internal/service/*_test.go` — юнит-тесты сервисного слоя.

**Принцип:** тесты работают **без базы данных**. Вместо реального репозитория передаётся `mock`:

```go
// workout_test.go — мок репозитория
type mockWorkoutRepo struct {
    workouts map[int]*models.Workout
    nextID   int
}
func (m *mockWorkoutRepo) Create(_ context.Context, w *models.Workout) (int, error) {
    // просто кладём в map
}
```

**Запуск:**
```bash
go test ./internal/service/... -v   # 23 теста, все проходят
```

**Что проверяется:**
- `TestWorkoutGetByID_NotOwner` — чужая тренировка возвращает `ErrForbidden`
- `TestLogin_WrongPassword` — неверный пароль возвращает `ErrInvalidCredentials`
- `TestParseToken_WrongSecret` — JWT с другим секретом отклоняется

**Почему моки, а не реальная БД:** тесты должны быть быстрыми и не зависеть от инфраструктуры. Это стало возможным благодаря принципу D из SOLID — сервисы зависят от интерфейсов.

---

### Нагрузка (Bench)

**Где:** `internal/repository/bench_test.go` — бенчмарки репозиторного слоя.

```go
func BenchmarkExerciseRepository_Create(b *testing.B) {
    pool := benchPool(b)      // пропустит тест если нет БД
    repo := NewExerciseRepository(pool)
    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        repo.Create(ctx, &models.Exercise{...})
    }
}
```

**Запуск:**
```bash
TEST_DATABASE_URL=postgres://... go test ./internal/repository/... -bench=. -benchmem
```

**Результат показывает:**
- `ns/op` — наносекунд на операцию
- `B/op` — байт выделено на операцию
- `allocs/op` — количество аллокаций

Доступные бенчмарки: `BenchmarkExerciseRepository_Create`, `BenchmarkWorkoutRepository_Create`, `BenchmarkUserRepository_GetByID`, `BenchmarkWorkoutRepository_ListByUser`.

---

### Подключения к БД

**Где:** `internal/repository/postgres.go`:

```go
pool, err := pgxpool.NewWithConfig(ctx, cfg)
pool.Ping(ctx)  // проверяем что БД доступна при старте
```

**Что такое пул соединений:**
Открыть новое TCP-соединение к PostgreSQL занимает ~5-10 мс. При 100 запросах в секунду это 0.5-1 с только на соединения. Пул держит несколько соединений открытыми и раздаёт их на запросы.

**pgxpool.Pool** — потокобезопасен: несколько горутин могут одновременно брать соединения из пула. Мы создаём пул **один раз** в `cmd/main.go` и передаём всем репозиториям по указателю.

**DATABASE_URL формат:**
```
postgres://postgres:postgres@localhost:5433/workout_tracker?sslmode=disable
```

---

### Публикация пакетов

**go.mod** — файл модуля. Определяет имя модуля и его зависимости:

```
module github.com/churilovmn1/workout-tracker   ← уникальное имя
go 1.25.0

require (
    github.com/go-chi/chi/v5 v5.3.0
    github.com/golang-jwt/jwt/v5 v5.3.1
    ...
)
```

**Как Go находит пакеты:** по URL из `module` — это путь на GitHub. `go get github.com/go-chi/chi/v5` скачивает пакет, добавляет в `go.sum` контрольные суммы.

**go mod tidy** — убирает неиспользуемые зависимости, добавляет недостающие. Мы запускали после удаления Telegram-бота.

---

### Не больше 5 параметров метода

**Принцип:** если у функции больше 5 параметров — это сигнал что она делает слишком много или нужна структура.

**В проекте** большинство методов имеют 1–3 параметра. Исключение — `NewRouter` (7 аргументов), где все параметры — сервисы всего приложения. Это нормально для точки сборки.

**Паттерн для решения** — структуры как конфиг:
```go
// Вместо:
func New(a, b, c, d, e, f string) {}
// Лучше:
type Config struct { A, B, C, D, E, F string }
func New(cfg Config) {}
```

---

### ...string (variadic / вариативные параметры)

**Где в проекте:** `internal/repository/exercise.go`, функция `WithSearch`:

```go
func WithSearch(term string, fields ...string) ExerciseFilterOption {
    return func(f *exerciseFilter) {
        // ...
        for i, field := range fields {
            parts[i] = field + " ILIKE " + placeholder
        }
    }
}
```

**Использование:**
```go
WithSearch(search, "name", "description")
// fields = ["name", "description"]
// Строит: name ILIKE $1 OR description ILIKE $1
```

`...string` позволяет передать любое количество строк. Внутри функции это обычный `[]string`.

---

### Системные вызовы (syscall)

**Где:** `cmd/main.go`:

```go
signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
```

`SIGINT` — сигнал от Ctrl+C.
`SIGTERM` — сигнал завершения от ОС (например, `docker stop` посылает SIGTERM).

Это системный вызов к ядру ОС: «уведоми меня когда придёт этот сигнал». Go оборачивает низкоуровневый `kill(2)` syscall в удобный пакет `os/signal`.

---

### Безопасность

В проекте реализовано несколько уровней защиты:

| Угроза | Защита | Где в коде |
|--------|--------|-----------|
| Перехват пароля | bcrypt-хэш, передаётся только JWT | `service/auth.go` |
| Подделка JWT | Подпись HMAC-SHA256 | `service/auth.go` → `ParseToken` |
| SQL-инъекция | Параметризованные запросы `$1, $2` | Все файлы `repository/` |
| Доступ к чужим данным | Проверка `UserID` в сервисе | `service/workout.go` → `GetByID`, `Update` |
| Брутфорс | Rate limiter 100 req/min на IP | `handler/ratelimit.go` |
| Неавторизованный доступ | `AuthMiddleware` на всех `/api/*` | `handler/middleware.go` |
| Доступ к admin-функциям | `AdminOnly` middleware | `handler/middleware.go` |
| Раскрытие внутренностей | `/debug/pprof` только для admin | `handler/router.go` |

---

### uintptr not safe

`uintptr` — целое число, содержащее адрес в памяти. **Не является ссылкой** с точки зрения сборщика мусора.

**Проблема:**
```go
ptr := uintptr(unsafe.Pointer(&obj))
// GC может переместить obj в памяти прямо здесь
// ptr теперь указывает на мусор
```

**Почему это опасно:** GC Go перемещает объекты при сборке мусора (compacting GC). `unsafe.Pointer` обновляется вместе с объектом, а `uintptr` — нет.

**В нашем проекте** типобезопасные ключи context:
```go
// handler/middleware.go
type contextKey string  // не uintptr — компилятор проверяет тип

const ctxUserID contextKey = "user_id"

ctx := context.WithValue(r.Context(), ctxUserID, claims.UserID)
id, _ := r.Context().Value(ctxUserID).(int)  // type assertion
```

---

## 4. Быстрые ответы на вопросы преподавателя

**— Зачем три слоя (handler/service/repository)?**
Разделение ответственности. Можно заменить PostgreSQL на другую БД, не трогая бизнес-логику. Можно добавить CLI, не переписывая HTTP-обработчики. И тесты можно писать без базы данных.

**— Как работает JWT?**
При входе сервер создаёт токен: JSON-заголовок + JSON-payload с user_id и role, подписанный HMAC-SHA256. Клиент хранит токен и отправляет в каждом запросе. Сервер проверяет подпись — не нужно обращаться к БД при каждом запросе.

**— Почему пул соединений, а не одно соединение?**
HTTP-сервер обрабатывает несколько запросов одновременно (каждый в горутине). Одно соединение к БД — бутылочное горлышко, запросы будут ждать в очереди. Пул даёт каждой горутине своё соединение.

**— Что такое транзакция и зачем она в Create?**
Транзакция — набор операций, которые выполняются целиком или не выполняются вовсе. При создании тренировки: если вставился `workout`, но упал INSERT одного из `workout_exercises` — без транзакции в БД останется «голая» тренировка без упражнений. С транзакцией — rollback, всё откатится.

**— Как тесты работают без БД?**
Через мок-объекты. Интерфейс `workoutRepository` определяет что нужно сервису. В тесте реализуем этот интерфейс через обычную map в памяти. Компилятор проверяет что мок реализует все методы.

**— Что такое middleware?**
Функция-обёртка над handler-ом. `AuthMiddleware` проверяет JWT до того, как запрос дойдёт до handler-а. Если токен невалиден — возвращает 401 сразу. Handler даже не вызывается.
