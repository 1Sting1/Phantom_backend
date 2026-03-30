# Детальная схема проекта веб-приложения для мессенджера

## 1. ОБЩАЯ АРХИТЕКТУРА МИКРОСЕРВИСОВ

### Основные микросервисы:

1. **API Gateway Service**
   - Единая точка входа
   - Маршрутизация запросов
   - Rate limiting
   - JWT валидация

2. **Auth Service**
   - Регистрация/авторизация
   - Управление сессиями
   - OAuth2 интеграция

3. **User Service**
   - Профили пользователей
   - Личный кабинет
   - Настройки

4. **Forum Service**
   - Система форума
   - Темы, посты, комментарии
   - Модерация

6. **Sticker Service?**
   - Каталог стикеров
   - Скачанные стикеры

8. **Localization Service**
   - Управление переводами
   - Кэширование языковых пакетов

9. **Notification Service**
   - Email уведомления
   - Push-уведомления
   - Внутренние уведомления

10. **File Storage Service**
    - Загрузка/хранение файлов
    - Аватары, стикеры
    - CDN интеграция

---

## 2. СТРУКТУРА КАЖДОГО МИКРОСЕРВИСА

### Стандартная структура Go микросервиса:

```
service-name/
├── cmd/
│   └── main.go
├── internal/
│   ├── handlers/
│   ├── services/
│   ├── repository/
│   ├── models/
│   ├── middleware/
│   └── config/
├── pkg/
│   ├── utils/
│   ├── logger/
│   └── validator/
├── migrations/
├── config/
├── tests/
├── proto/ (для gRPC)
├── docker/
├── go.mod
├── go.sum
└── README.md
```

---

## 3. ДЕТАЛЬНАЯ СХЕМА МИКРОСЕРВИСОВ

### 3.1 API GATEWAY SERVICE

**Директории и файлы:**

**cmd/**
- `main.go` - точка входа приложения

**internal/handlers/**
- `router.go` - настройка маршрутизации
- `proxy_handler.go` - проксирование запросов к сервисам
- `health_handler.go` - health check endpoints

**internal/middleware/**
- `auth_middleware.go` - проверка JWT токенов
- `rate_limiter.go` - ограничение запросов
- `cors_middleware.go` - настройка CORS
- `logger_middleware.go` - логирование запросов
- `recovery_middleware.go` - обработка паник

**internal/config/**
- `config.go` - структура конфигурации
- `loader.go` - загрузка конфигурации из env/файлов

**internal/services/**
- `service_registry.go` - реестр доступных сервисов
- `load_balancer.go` - балансировка нагрузки

**pkg/logger/**
- `logger.go` - настройка логирования (zap/logrus)

**config/**
- `config.yaml` - конфигурационный файл
- `services.yaml` - список микросервисов и их адреса

---

### 3.2 AUTH SERVICE

**internal/handlers/**
- `auth_handler.go` - регистрация, логин, logout
- `token_handler.go` - refresh token, валидация
- `oauth_handler.go` - OAuth2 провайдеры (Google, GitHub)
- `password_handler.go` - сброс/восстановление пароля

**internal/services/**
- `auth_service.go` - бизнес-логика аутентификации
- `token_service.go` - генерация/валидация JWT
- `oauth_service.go` - интеграция с OAuth провайдерами
- `email_service.go` - отправка писем подтверждения

**internal/repository/**
- `user_repository.go` - работа с таблицей users
- `session_repository.go` - управление сессиями
- `token_repository.go` - хранение refresh токенов

**internal/models/**
- `user.go` - модель пользователя
- `session.go` - модель сессии
- `token.go` - модель токена

**internal/middleware/**
- `validation_middleware.go` - валидация входящих данных

**migrations/**
- `001_create_users_table.up.sql`
- `001_create_users_table.down.sql`
- `002_create_sessions_table.up.sql`
- `002_create_sessions_table.down.sql`
- `003_create_refresh_tokens_table.up.sql`
- `003_create_refresh_tokens_table.down.sql`

**Database Tables:**
- users (id, email, password_hash, created_at, updated_at, email_verified, is_active)
- sessions (id, user_id, token, ip_address, user_agent, expires_at, created_at)
- refresh_tokens (id, user_id, token, expires_at, created_at)
- oauth_connections (id, user_id, provider, provider_user_id, created_at)

---

### 3.3 USER SERVICE

**internal/handlers/**
- `profile_handler.go` - получение/обновление профиля
- `settings_handler.go` - настройки пользователя
- `avatar_handler.go` - загрузка аватара
- `preferences_handler.go` - пользовательские предпочтения

**internal/services/**
- `profile_service.go` - управление профилем
- `settings_service.go` - работа с настройками
- `avatar_service.go` - обработка аватаров
- `statistics_service.go` - статистика пользователя

**internal/repository/**
- `profile_repository.go` - работа с профилями
- `settings_repository.go` - хранение настроек
- `statistics_repository.go` - сбор статистики

**internal/models/**
- `profile.go` - расширенный профиль
- `settings.go` - настройки
- `statistics.go` - статистика

**migrations/**
- `001_create_profiles_table.up.sql`
- `002_create_settings_table.up.sql`
- `003_create_user_statistics_table.up.sql`

**Database Tables:**
- profiles (id, user_id, display_name, bio, avatar_url, country, timezone, created_at, updated_at)
- user_settings (id, user_id, language, theme, notifications_enabled, privacy_level)
- user_statistics (id, user_id, messages_sent, forum_posts, stickers_purchased, last_activity)

---

### 3.4 FORUM SERVICE

**Предложение по реализации форума:**

**Структура форума:**
- Категории (Installation Issues, Feature Requests, General Discussion, Development)
- Подкатегории
- Темы (threads)
- Посты
- Комментарии/ответы
- Система лайков/дизлайков
- Система репутации пользователей
- Теги
- Поиск по форуму
- Модерация (бан, удаление, закрепление тем)
- Подписки на темы
- Уведомления о новых ответах

**internal/handlers/**
- `category_handler.go` - управление категориями
- `thread_handler.go` - создание/редактирование тем
- `post_handler.go` - посты и комментарии
- `reaction_handler.go` - лайки/реакции
- `search_handler.go` - поиск по форуму
- `moderation_handler.go` - модераторские функции
- `subscription_handler.go` - подписки на темы

**internal/services/**
- `category_service.go` - бизнес-логика категорий
- `thread_service.go` - управление темами
- `post_service.go` - работа с постами
- `reaction_service.go` - система реакций
- `search_service.go` - поисковый движок
- `moderation_service.go` - модерация контента
- `reputation_service.go` - система репутации
- `notification_service.go` - уведомления форума

**internal/repository/**
- `category_repository.go`
- `thread_repository.go`
- `post_repository.go`
- `reaction_repository.go`
- `tag_repository.go`
- `subscription_repository.go`
- `moderation_repository.go`

**internal/models/**
- `category.go`
- `thread.go`
- `post.go`
- `reaction.go`
- `tag.go`
- `subscription.go`

**migrations/**
- `001_create_categories_table.up.sql`
- `002_create_threads_table.up.sql`
- `003_create_posts_table.up.sql`
- `004_create_reactions_table.up.sql`
- `005_create_tags_table.up.sql`
- `006_create_thread_tags_table.up.sql`
- `007_create_subscriptions_table.up.sql`
- `008_create_moderator_actions_table.up.sql`

**Database Tables:**
- categories (id, name, description, slug, parent_id, order, icon, created_at)
- threads (id, category_id, user_id, title, slug, content, is_pinned, is_locked, views_count, created_at, updated_at)
- posts (id, thread_id, user_id, parent_post_id, content, is_edited, edited_at, created_at)
- reactions (id, post_id, user_id, reaction_type, created_at)
- tags (id, name, slug, color)
- thread_tags (thread_id, tag_id)
- subscriptions (id, user_id, thread_id, created_at)
- moderator_actions (id, moderator_id, action_type, target_type, target_id, reason, created_at)
- user_reputation (id, user_id, reputation_score, helpful_posts, total_posts)

---

### 3.6 STICKER SERVICE

**internal/handlers/**
- `sticker_handler.go` - каталог стикеров
- `pack_handler.go` - наборы стикеров
- `order_handler.go` - скачанные стикеры
- `category_handler.go` - категории стикеров

**internal/services/**
- `sticker_service.go` - управление стикерами
- `cart_service.go` - логика корзины
- `inventory_service.go` - управление запасами

**internal/repository/**
- `sticker_repository.go`
- `pack_repository.go`
- `cart_repository.go`
- `order_repository.go`
- `category_repository.go`

**internal/models/**
- `sticker.go`
- `pack.go`
- `cart.go`
- `order.go`
- `category.go`

**migrations/**
- `001_create_sticker_categories_table.up.sql`
- `002_create_stickers_table.up.sql`
- `003_create_sticker_packs_table.up.sql`
- `004_create_carts_table.up.sql`
- `005_create_cart_items_table.up.sql`
- `006_create_orders_table.up.sql`
- `007_create_order_items_table.up.sql`

**Database Tables:**
- sticker_categories (id, name, slug, description, icon)
- stickers (id, name, description, preview_url, file_url, price, category_id, created_at, is_active)
- sticker_packs (id, name, description, preview_url, price, discount, created_at)
- pack_stickers (pack_id, sticker_id)
- carts (id, user_id, created_at, updated_at)
- cart_items (id, cart_id, sticker_id, pack_id, quantity)
- orders (id, user_id, total_amount, status, payment_id, created_at, completed_at)
- order_items (id, order_id, sticker_id, pack_id, price, quantity)
- user_stickers (id, user_id, sticker_id, purchased_at)

---

### 3.8 LOCALIZATION SERVICE

**internal/handlers/**
- `translation_handler.go` - получение переводов
- `language_handler.go` - список языков
- `admin_handler.go` - управление переводами (для админов)

**internal/services/**
- `translation_service.go` - загрузка переводов
- `cache_service.go` - кэширование языковых пакетов
- `fallback_service.go` - фоллбэк на английский

**internal/repository/**
- `translation_repository.go`
- `language_repository.go`

**internal/models/**
- `translation.go`
- `language.go`

**migrations:**
- `001_create_languages_table.up.sql`
- `002_create_translations_table.up.sql`

**Database Tables:**
- languages (id, code, name, native_name, is_active, is_rtl)
- translations (id, language_id, key, value, context, created_at, updated_at)

---

### 3.9 NOTIFICATION SERVICE

**internal/handlers/**
- `notification_handler.go` - получение уведомлений
- `preference_handler.go` - настройки уведомлений
- `mark_read_handler.go` - пометка прочитанными

**internal/services/**
- `notification_service.go` - создание уведомлений
- `email_service.go` - отправка email
- `push_service.go` - push-уведомления
- `websocket_service.go` - real-time уведомления
- `template_service.go` - шаблоны уведомлений

**internal/repository/**
- `notification_repository.go`
- `preference_repository.go`
- `template_repository.go`

**internal/models/**
- `notification.go`
- `preference.go`
- `template.go`

**migrations:**
- `001_create_notifications_table.up.sql`
- `002_create_notification_preferences_table.up.sql`
- `003_create_notification_templates_table.up.sql`

**Database Tables:**
- notifications (id, user_id, type, title, message, link, is_read, created_at)
- notification_preferences (id, user_id, email_enabled, push_enabled, forum_replies, new_stickers, system_updates)
- notification_templates (id, type, language_id, subject, body)

---

### 3.10 FILE STORAGE SERVICE

**internal/handlers/**
- `upload_handler.go` - загрузка файлов
- `download_handler.go` - скачивание файлов
- `delete_handler.go` - удаление файлов

**internal/services/**
- `storage_service.go` - работа с хранилищем
- `s3_service.go` - интеграция с S3
- `image_processor_service.go` - обработка изображений
- `validation_service.go` - валидация файлов

**internal/repository/**
- `file_repository.go` - метаданные файлов

**internal/models:**
- `file.go`

**migrations:**
- `001_create_files_table.up.sql`

**Database Tables:**
- files (id, user_id, filename, original_filename, size, mime_type, storage_path, url, created_at)

---

## 4. ОБЩИЕ КОМПОНЕНТЫ ДЛЯ ВСЕХ СЕРВИСОВ

### pkg/ (общие пакеты)

**pkg/database/**
- `postgres.go` - подключение к PostgreSQL
- `migrations.go` - запуск миграций
- `transaction.go` - работа с транзакциями

**pkg/redis/**
- `client.go` - подключение к Redis
- `cache.go` - кэширование

**pkg/jwt/**
- `token.go` - генерация/валидация JWT
- `claims.go` - структура claims

**pkg/validator/**
- `validator.go` - валидация структур
- `custom_rules.go` - кастомные правила

**pkg/logger/**
- `logger.go` - настройка логгера
- `fields.go` - структурированные поля

**pkg/errors/**
- `errors.go` - кастомные ошибки
- `codes.go` - коды ошибок
- `handler.go` - обработчик ошибок

**pkg/http/**
- `response.go` - стандартные HTTP ответы
- `pagination.go` - пагинация

**pkg/grpc/** (если используется gRPC между сервисами)
- `client.go` - gRPC клиент
- `server.go` - gRPC сервер
- `interceptors.go` - интерсепторы

---

## 5. ИНФРАСТРУКТУРА

### Дополнительные сервисы:

**Message Broker (RabbitMQ/Kafka)**
- Асинхронная коммуникация между сервисами
- Очереди задач
- Event-driven архитектура

**Redis**
- Кэширование
- Сессии
- Rate limiting
- Real-time данные

**Elasticsearch** (опционально)
- Полнотекстовый поиск по форуму
- Поиск по документации
- Аналитика

**Monitoring Stack**
- Prometheus - сбор метрик
- Grafana - визуализация
- Jaeger - трейсинг

---

## 6. КОНФИГУРАЦИОННЫЕ ФАЙЛЫ

### docker-compose.yml
- Описание всех сервисов
- PostgreSQL
- Redis
- RabbitMQ
- Traefik (для проксирования)

### Kubernetes манифесты (для продакшена)
- Deployments для каждого сервиса
- Services
- Ingress
- ConfigMaps
- Secrets

### CI/CD (.gitlab-ci.yml / .github/workflows)
- Тесты
- Сборка Docker образов
- Деплой

---

## 7. API ENDPOINTS СХЕМА

### API Gateway маршруты:

**Auth:**
- POST /api/v1/auth/register
- POST /api/v1/auth/login
- POST /api/v1/auth/logout
- POST /api/v1/auth/refresh
- POST /api/v1/auth/forgot-password
- POST /api/v1/auth/reset-password

**User:**
- GET /api/v1/user/profile
- PUT /api/v1/user/profile
- GET /api/v1/user/settings
- PUT /api/v1/user/settings
- POST /api/v1/user/avatar

**Forum:**
- GET /api/v1/forum/categories
- GET /api/v1/forum/threads
- POST /api/v1/forum/threads
- GET /api/v1/forum/threads/:id
- POST /api/v1/forum/posts
- POST /api/v1/forum/reactions
- GET /api/v1/forum/search

**Documentation:**
- GET /api/v1/docs
- GET /api/v1/docs/:slug
- GET /api/v1/docs/search
- GET /api/v1/docs/versions

**Shop:**
- GET /api/v1/shop/stickers
- GET /api/v1/shop/packs
- POST /api/v1/shop/order

**Localization:**
- GET /api/v1/localization/:language
- GET /api/v1/localization/languages

**Notifications:**
- GET /api/v1/notifications
- PUT /api/v1/notifications/:id/read
- GET /api/v1/notifications/preferences

---

## 8. МЕЖСЕРВИСНАЯ КОММУНИКАЦИЯ

- Websocket

---

## 9. БАЗА ДАННЫХ

- auth_db
- user_db
- forum_db
- documentation_db
- shop_db
- payment_db
- notification_db

### Общие таблицы:
- audit_logs (логирование всех действий)
- system_settings (глобальные настройки)

---

## 10. БЕЗОПАСНОСТЬ

### Компоненты безопасности в каждом сервисе:

**internal/middleware/**
- `rate_limiter.go` - защита от DDoS
- `input_sanitizer.go` - очистка входных данных
- `sql_injection_guard.go` - защита от SQL injection

**pkg/security/**
- `encryption.go` - шифрование данных
- `hash.go` - хеширование паролей (bcrypt)
- `csrf.go` - защита от CSRF

---

## 11. ТЕСТИРОВАНИЕ

### В каждом сервисе:

**tests/**
- `unit/` - юнит-тесты
  - handlers_test.go
  - services_test.go
  - repository_test.go
- `integration/` - интеграционные тесты
  - api_test.go
  - database_test.go
- `e2e/` - end-to-end тесты

---