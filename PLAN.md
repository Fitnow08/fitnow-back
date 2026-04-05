# Plan: Train Feature — полная реализация

## Статус реализации

| Компонент | Статус |
|-----------|--------|
| `internal/models/domain/train.go` | ✅ Создан |
| `internal/feature/train/entity.go` | ✅ Создан |
| `internal/feature/train/repository.go` | ✅ Создан |
| `internal/feature/train/service.go` | ✅ Создан |
| `internal/feature/train/handlers.go` | ✅ Создан |
| `internal/app/repositories.go` | ✅ Обновлён |
| `internal/app/services.go` | ✅ Обновлён |
| `internal/app/handlers.go` | ✅ Обновлён |
| `internal/handlers/handlers.go` | ⚠️ Частично (маршруты не полные) |

---

## Контекст

- БД: PostgreSQL, миграция `20260404164721_train.sql` создаёт таблицы `trains`, `train_exercises`, `exercises`, `user_trains`
- Авторизация: JWT через `internal/feature/auth/jwt.go` + middleware `internal/handlers/middleware/sso.go`
- Архитектура: handler → service → repository, по образцу `internal/feature/auth/`

---

## Эндпоинты

### Публичные (без JWT)

| Method | Route | Handler |
|--------|-------|---------|
| GET | `/api/v1/trains` | `TrainHandler.GetAllTrains` |
| GET | `/api/v1/trains/{id}` | `TrainHandler.GetTrainByID` |
| GET | `/api/v1/exercises` | `TrainHandler.GetAllExercises` |

### Защищённые (требуют `Authorization: Bearer <token>`)

| Method | Route | Handler |
|--------|-------|---------|
| POST | `/api/v1/trains` | `TrainHandler.CreateTrain` |
| PUT | `/api/v1/trains/{id}` | `TrainHandler.UpdateTrain` |
| DELETE | `/api/v1/trains/{id}` | `TrainHandler.DeleteTrain` |
| GET | `/api/v1/user/trains` | `TrainHandler.GetUserTrains` |
| POST | `/api/v1/user/trains/{id}` | `TrainHandler.AddUserTrain` |
| DELETE | `/api/v1/user/trains/{id}` | `TrainHandler.RemoveUserTrain` |
| POST | `/api/v1/exercises` | `TrainHandler.CreateExercise` |

---

## Целевая конфигурация маршрутов (`handlers.go`)

```go
r.Route("/trains", func(r chi.Router) {
    // Публичные
    r.Get("/", handlers.TrainHandler.GetAllTrains)
    r.Get("/{id}", handlers.TrainHandler.GetTrainByID)

    // Защищённые
    r.Group(func(r chi.Router) {
        r.Use(customMiddleware.AuthMiddleware(""))
        r.Post("/", handlers.TrainHandler.CreateTrain)
        r.Put("/{id}", handlers.TrainHandler.UpdateTrain)
        r.Delete("/{id}", handlers.TrainHandler.DeleteTrain)
    })
})

r.Route("/exercises", func(r chi.Router) {
    r.Get("/", handlers.TrainHandler.GetAllExercises)
    r.Group(func(r chi.Router) {
        r.Use(customMiddleware.AuthMiddleware(""))
        r.Post("/", handlers.TrainHandler.CreateExercise)
    })
})

r.Group(func(r chi.Router) {
    r.Use(customMiddleware.AuthMiddleware(""))
    r.Route("/user/trains", func(r chi.Router) {
        r.Get("/", handlers.TrainHandler.GetUserTrains)
        r.Post("/{id}", handlers.TrainHandler.AddUserTrain)
        r.Delete("/{id}", handlers.TrainHandler.RemoveUserTrain)
    })
})
```

---

## Текущее состояние маршрутов (`handlers.go`)

Зарегистрировано:
- `GET /api/v1/train/` → `GetAllTrains`
- `POST /api/v1/train/` → `CreateTrain` (под AuthMiddleware)
- `PUT /api/v1/train/` → `UpdateTrain` (под AuthMiddleware, без `{id}`)
- `DELETE /api/v1/train/` → `DeleteTrain` (под AuthMiddleware, без `{id}`)

Не зарегистрировано:
- `GET /trains/{id}` → `GetTrainByID`
- `GET /exercises` → `GetAllExercises`
- `POST /exercises` → `CreateExercise`
- `GET /user/trains` → `GetUserTrains`
- `POST /user/trains/{id}` → `AddUserTrain`
- `DELETE /user/trains/{id}` → `RemoveUserTrain`

---

## Domain модели (`internal/models/domain/train.go`)

```go
type Train struct {
    ID         uuid.UUID `json:"id"`
    Title      string    `json:"title"`
    Type       string    `json:"type"`
    Duration   int64     `json:"duration"`
    IsPublic   bool      `json:"is_public"`
    Difficulty string    `json:"difficulty"`
    Calories   int64     `json:"calories"`
    CreatedBy  uuid.UUID `json:"created_by"`
    CreatedAt  time.Time `json:"created_at"`
}

type Exercise struct {
    ID          uuid.UUID `json:"id"`
    Title       string    `json:"title"`
    Description string    `json:"description"`
}
```

---

## DTO / Валидация (`internal/feature/train/entity.go`)

```go
type CreateTrainRequest struct {
    Title      string `json:"title" validate:"required"`
    Type       string `json:"type" validate:"required,oneof=strength cardio stretching"`
    Duration   int64  `json:"duration" validate:"required,min=1"`
    IsPublic   bool   `json:"is_public"`
    Difficulty string `json:"difficulty" validate:"required,oneof=easy medium hard"`
    Calories   int64  `json:"calories" validate:"required,min=0"`
}

type UpdateTrainRequest struct {
    Title      string `json:"title"`
    Type       string `json:"type" validate:"omitempty,oneof=strength cardio stretching"`
    Duration   int64  `json:"duration" validate:"omitempty,min=1"`
    IsPublic   *bool  `json:"is_public"`
    Difficulty string `json:"difficulty" validate:"omitempty,oneof=easy medium hard"`
    Calories   int64  `json:"calories" validate:"omitempty,min=0"`
}

type CreateExerciseRequest struct {
    Title       string `json:"title" validate:"required"`
    Description string `json:"description"`
}
```

---

## Паттерн получения user_id в защищённых хендлерах

```go
claims := r.Context().Value(constants.UserClaimsKey).(*constants.UserClaims)
userID := claims.ID
```

---

## Таблицы БД (из миграции)

```sql
CREATE TYPE train_type AS ENUM ('strength', 'cardio', 'stretching');
CREATE TYPE difficulty_level AS ENUM ('easy', 'medium', 'hard');

CREATE TABLE exercises (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    description TEXT
);

CREATE TABLE trains (
    id UUID DEFAULT uuid_generate_v4() PRIMARY KEY,
    title TEXT NOT NULL,
    type train_type,
    duration BIGINT,
    is_public BOOLEAN DEFAULT true NOT NULL,
    difficulty difficulty_level,
    created_by UUID REFERENCES users(id),
    calories BIGINT NOT NULL,
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    version BIGINT NOT NULL DEFAULT 0
);

CREATE TABLE train_exercises (
    id UUID DEFAULT uuid_generate_v4() NOT NULL,
    steps INT,
    sets INT,
    position INT NOT NULL,
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    exercises_id UUID NOT NULL REFERENCES exercises(id),
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (id)
);

CREATE TABLE user_trains (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    train_id UUID NOT NULL REFERENCES trains(id) ON DELETE CASCADE,
    added_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,
    PRIMARY KEY (user_id, train_id)
);
```

---

## Верификация

```bash
# Компиляция
go build ./...

# Публичные
curl GET /api/v1/trains             # список публичных тренировок
curl GET /api/v1/trains/{id}        # одна тренировка
curl GET /api/v1/exercises          # список упражнений

# Защищённые (нужен JWT)
curl -H "Authorization: Bearer <token>" POST /api/v1/trains           # создать
curl -H "Authorization: Bearer <token>" PUT /api/v1/trains/{id}       # обновить (только создатель)
curl -H "Authorization: Bearer <token>" DELETE /api/v1/trains/{id}    # удалить (только создатель)
curl -H "Authorization: Bearer <token>" GET /api/v1/user/trains       # мои тренировки
curl -H "Authorization: Bearer <token>" POST /api/v1/user/trains/{id} # добавить к себе
curl -H "Authorization: Bearer <token>" DELETE /api/v1/user/trains/{id} # убрать
curl -H "Authorization: Bearer <token>" POST /api/v1/exercises        # создать упражнение
```

---

## Следующий шаг

Обновить `internal/handlers/handlers.go` — заменить текущий блок `/train` на полную конфигурацию маршрутов (см. раздел "Целевая конфигурация маршрутов" выше).
