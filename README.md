# User Segment Service

Сервис для управления сегментами пользователей и проведения экспериментов (A/B-тестов) в рамках новых функций VK.

## 🚀 Описание

Сервис позволяет:

- Создавать, изменять и удалять сегменты пользователей
- Добавлять и удалять пользователей из сегментов
- Распределять сегменты случайным образом на заданный процент пользователей
- Получать список сегментов по `user_id` через API

---

## 🛠 Технологии

- Go 1.21
- Gin (REST API)
- GORM (Postgres ORM)
- PostgreSQL
- Docker / docker-compose

---

## ⚙️ Запуск проекта

1. **Клонировать репозиторий**

```bash
git clone https://github.com/yourusername/user-segment-service-pg.git
cd user-segment-service-pg
```
2. **Собрать и запустить через docker-compose**

```bash
docker-compose up --build
```
3. **Запустить программу**
```bash
go run cmd/main.go
```
## 📖 Основные эндпоинты
| Метод  | URL                           | Описание                                         |
| ------ | ----------------------------- | ------------------------------------------------ |
| GET    | `/users/:id/segments`         | Получить список сегментов для пользователя       |
| POST   | `/segments`                   | Создать новый сегмент                            |
| PUT    | `/segments/:id`               | Обновить сегмент                                 |
| DELETE | `/segments/:id`               | Удалить сегмент                                  |
| POST   | `/segments/:id/users/:userId` | Добавить пользователя в сегмент                  |
| DELETE | `/segments/:id/users/:userId` | Удалить пользователя из сегмента                 |
| POST   | `/segments/:id/distribute`    | Распределить сегмент случайно на % пользователей |


## 📝 Примеры запросов

### Создать сегмент

```bash
curl -X POST http://localhost:8080/segments \
  -H "Content-Type: application/json" \
  -d '{"name":"MAIL_GPT"}'
```
### Получить сегменты пользователя
```bash
curl http://localhost:8080/users/1/segments
```
