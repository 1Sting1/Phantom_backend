# Phantom Backend

Микросервисная архитектура для мессенджера Phantom.

## Структура проекта

- `services/` - микросервисы
  - `api-gateway/` - API Gateway (порт 8080)
  - `auth-service/` - Сервис аутентификации (порт 8001)
  - `user-service/` - Сервис пользователей (порт 8002)
  - `forum-service/` - Сервис форума (порт 8003)
  - `sticker-service/` - Сервис стикеров (порт 8004)
  - `localization-service/` - Сервис локализации (порт 8005)
  - `notification-service/` - Сервис уведомлений (порт 8006)
- `pkg/` - общие пакеты
- `frontend/` - React frontend

## API Endpoints

### Auth
- POST `/api/v1/auth/register` - Регистрация (body: `email`, `password`). Ответ: `token`, `refresh_token`, `user`.
- POST `/api/v1/auth/login` - Вход (body: `email`, `password`). Ответ: `token`, `refresh_token`, `user`.
- POST `/api/v1/auth/logout` - Выход (опционально body: `refresh_token` — инвалидирует refresh-токен).
- POST `/api/v1/auth/refresh` - Обновление токена (body: `refresh_token`). Ответ: новый `token`, `refresh_token`, `user` (запоминание пользователя на 30 дней).
- GET `/api/v1/me` - Текущий пользователь (заголовок `Authorization: Bearer <token>`).

### User
- GET `/api/v1/user/profile` - Получить профиль (в т.ч. `display_name` — никнейм)
- PUT `/api/v1/user/profile` - Обновить профиль (body: `display_name`, `bio`, `avatar_url`, …)
- POST `/api/v1/user/avatar` - Установить аватар (body: `avatar_url`)
- GET `/api/v1/user/settings` - Получить настройки
- PUT `/api/v1/user/settings` - Обновить настройки

### Forum
- GET `/api/v1/forum/categories` - Список категорий (дефолтные: general, installation, features, development)
- GET `/api/v1/forum/threads` - Список тем
- POST `/api/v1/forum/threads` - Создать тему (требуется авторизация; body: `title`, `content`, `category_id`)
- GET `/api/v1/forum/threads/:id` - Получить тему
- POST `/api/v1/forum/posts` - Создать пост (требуется авторизация; body: `thread_id`, `content`)

### Shop
- GET `/api/v1/shop/stickers` - Список стикеров
- GET `/api/v1/shop/packs` - Список наборов
- POST `/api/v1/shop/order` - Создать заказ

### Localization
- GET `/api/v1/localization/:language` - Получить переводы
- GET `/api/v1/localization/languages` - Список языков

### Notifications
- GET `/api/v1/notifications` - Получить уведомления
- PUT `/api/v1/notifications/:id/read` - Отметить как прочитанное
- GET `/api/v1/notifications/preferences` - Настройки уведомлений

## Базы данных

Каждый сервис использует свою базу данных:
- `auth_db` - для auth-service
- `user_db` - для user-service
- `forum_db` - для forum-service
- `shop_db` - для sticker-service

## Миграции

`migrations/` каждого сервиса. Для применения миграций используйте `golang-migrate` или вручную.

- **auth-service**: после `001_create_users_table` выполните `002_create_refresh_tokens_table.up.sql` для работы refresh-токенов и «запомнить пользователя».
- **forum-service**: после `001_create_forum_tables` выполните `002_seed_default_categories.up.sql` для дефолтных категорий (нужны для создания тем).

## Разделение frontend/backend (Docker)

Чтобы запустить frontend и backend в разных репозиториях, используйте отдельные compose-файлы:

- Backend (Traefik + API-gateway + сервисы + БД): `docker-compose.backend.yml`
- Frontend (Next.js): `frontend/docker-compose.yml`

Важно: frontend должен подключаться к тому же docker-network, что и Traefik из backend-стека (`app-network`), поэтому backend запускайте первым.

Пример запуска из текущего каталога:

```bash
docker compose -f docker-compose.backend.yml up -d
docker compose -f frontend/docker-compose.yml up -d
```

