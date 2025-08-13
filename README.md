### TaskFlow — Kanban + Time Tracking для команд и соло-разработчиков
Удобное управление проектами и задачами с аналитикой времени. Учебный open-source проект: практикуем Go на бэкенде и Vue 3 на фронтенде.

---

### Ключевые возможности
- **Задачи и проекты**: CRUD, статусы, дедлайны, приоритеты, метки
- **Kanban-доска**: drag-n-drop между колонками
- **Тайм-трекер**: старт/пауза/стоп, суммарное время по задачам/проектам
- **Файлы**: вложения к задачам (MinIO/S3 или локально)
- **Команды**: приглашения, роли (админ/участник)
- **Real-time**: уведомления по WebSocket (изменения задач, дедлайны)
- **Аутентификация**: JWT, refresh токены

---

### Архитектура
- **Бэкенд (Go)**: REST API на Gin/Echo, PostgreSQL (через `pgx` или GORM), миграции, JWT, WebSocket (gorilla/websocket), файловое хранилище (MinIO/S3), слои `handler → service → repository`.
- **Фронтенд (Vue 3)**: SPA на Composition API, состояние в Pinia, маршрутизация в Vue Router, HTTP через Axios, стили Tailwind CSS/Quasar, графики на Chart.js, drag-n-drop (`vue-draggable-next`).
- **Инфраструктура**: Docker Compose (PostgreSQL, MinIO, локальные контейнеры backend/frontend), опционально NGINX как reverse-proxy.

Взаимодействия:
- Клиент вызывает REST `api/v1/*`
- События (создание/обновление задач, тик таймера) рассылаются подписчикам через WebSocket канал
- Файлы загружаются на MinIO; в БД хранятся метаданные ссылок

---

### Доменная модель (минимум)
- **User**: id, email, имя, роль, hash пароля
- **Project**: id, название, описание, владелец
- **Membership**: user_id, project_id, роль
- **Task**: id, project_id, заголовок, описание, статус, дедлайн, порядок, метки
- **TimeEntry**: id, task_id, user_id, начало, конец, длительность
- **Attachment**: id, task_id, author_id, file_key, mime, размер
- **Notification**: id, type, payload, прочитано

---

### Планируемая структура репозитория
```text
backend/
  cmd/api/main.go
  internal/
    http/
      handlers/
      middleware/
      router.go
    domain/
      user/
      project/
      task/
      timeentry/
      attachment/
    repository/
      postgres/
    service/
    auth/
    ws/
  migrations/
  go.mod

frontend/
  src/
    pages/
    components/
    stores/
    router/
    api/
    styles/
  vite.config.ts
  package.json

deploy/
  docker-compose.yml
  backend.Dockerfile
  frontend.Dockerfile
  nginx/nginx.conf

.github/workflows/ci.yml
```

---

### API (обзор)
- Auth: `POST /api/v1/auth/register`, `POST /api/v1/auth/login`, `POST /api/v1/auth/refresh`
- Projects: `GET/POST /api/v1/projects`, `GET/PATCH/DELETE /api/v1/projects/:id`
- Tasks: `GET/POST /api/v1/projects/:id/tasks`, `GET/PATCH/DELETE /api/v1/tasks/:id`, `POST /api/v1/tasks/:id/reorder`
- Time tracking: `POST /api/v1/tasks/:id/time/start|pause|stop`, `GET /api/v1/tasks/:id/time`
- Attachments: `POST /api/v1/tasks/:id/attachments` (multipart), `GET /api/v1/attachments/:id`
- WebSocket: `GET /ws` (авторизация по JWT в query/header)

Пример: создание задачи
```http
POST /api/v1/projects/42/tasks
Content-Type: application/json
Authorization: Bearer <token>

{
  "title": "Сверстать доску",
  "description": "Колонки To Do / In Progress / Done",
  "status": "todo",
  "labels": ["ui", "p1"],
  "dueDate": "2025-08-20"
}
```

События WebSocket (тип → payload):
- `task.created`, `task.updated`, `task.moved`, `task.deleted`
- `time.started`, `time.paused`, `time.stopped`
- `attachment.added`

---

### Быстрый старт (локально)
Требования: Go ≥ 1.22, Node.js ≥ 20, Docker + Docker Compose.

1) Склонировать репозиторий и создать `.env` файлы (см. ниже).
2) Поднять инфраструктуру (PostgreSQL, MinIO):
```bash
docker compose -f deploy/docker-compose.yml up -d
```
3) Бэкенд (из `backend/`):
```bash
go mod tidy
go run ./cmd/api
```
4) Фронтенд (из `frontend/`):
```bash
npm i
npm run dev
```

Пример `deploy/docker-compose.yml` (набросок для реализации):
```yaml
version: "3.9"
services:
  db:
    image: postgres:16
    environment:
      POSTGRES_DB: taskflow
      POSTGRES_USER: taskflow
      POSTGRES_PASSWORD: taskflow
    ports: ["5432:5432"]
    volumes: ["pgdata:/var/lib/postgresql/data"]

  minio:
    image: minio/minio:latest
    command: server /data --console-address ":9001"
    environment:
      MINIO_ROOT_USER: minio
      MINIO_ROOT_PASSWORD: minio123
    ports: ["9000:9000", "9001:9001"]
    volumes: ["minio:/data"]

volumes:
  pgdata: {}
  minio: {}
```

Пример `backend/.env`:
```env
APP_ENV=dev
APP_PORT=8080
DB_DSN=postgres://taskflow:taskflow@localhost:5432/taskflow?sslmode=disable
JWT_SECRET=change-me
S3_ENDPOINT=http://localhost:9000
S3_ACCESS_KEY=minio
S3_SECRET_KEY=minio123
S3_BUCKET=taskflow
```

Пример `frontend/.env`:
```env
VITE_API_BASE=http://localhost:8080/api/v1
VITE_WS_URL=ws://localhost:8080/ws
```

---

### Разработка и тесты
- Бэкенд:
  - запуск: `go run ./cmd/api`
  - тесты: `go test ./...`
  - линт: `golangci-lint run` (рекомендуется)
- Фронтенд:
  - запуск: `npm run dev`
  - тесты: `vitest` или `jest` (на выбор)
  - линт: `eslint .` + `prettier --check .`

Рекомендуемые соглашения:
- **Conventional Commits** (`feat:`, `fix:`, `docs:`, `refactor:` ...)
- Ветки: `feature/<slug>`, `fix/<slug>`
- PR чек-лист: тесты проходят, описание изменений, скриншоты UI

---

### Дорожная карта
- M1: Базовые сущности (User/Project/Task), CRUD API, миграции
- M2: JWT аутентификация и роли, Pinia + маршрутизация
- M3: Тайм-трекер (старт/пауза/стоп), агрегации по времени
- M4: Kanban UI (перетаскивание), сортировка задач
- M5: Уведомления по WebSocket, live-обновления
- M6: Вложения (MinIO), предпросмотр
- M7: Дашборд аналитики (Chart.js), фильтры и отчеты
- M8: Интеграции (GitHub OAuth), экспорт CSV/PDF

---

### Лицензия
MIT. Используйте проект свободно для обучения и реальных задач.

---

### Зачем это полезно
- **Go**: продакшен-практики (REST, миграции, тесты, WebSocket, S3)
- **Vue 3**: SPA архитектура, управление состоянием, работа с real-time
- **Инфра**: Docker Compose, локальная разработка, CI/CD (можно добавить позже)

Если хотите — добавим CI, шаблоны issues/PR и автосборку контейнеров.
