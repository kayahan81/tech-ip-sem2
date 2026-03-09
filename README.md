---
# Практическое задание 17

## ЭФМО-02-25 

## Алиев Каяхан Командар оглы
---
# Тема работы
Разделение монолитного приложения на 2 микросервиса. Взаимодействие через HTTP

## Цели занятия
Научиться декомпозировать небольшую систему на два сервиса и организовать корректное синхронное взаимодействие по HTTP (с таймаутами, статусами ошибок и прокидыванием request-id).

## Структура проекта
<img width="426" height="752" alt="image" src="https://github.com/user-attachments/assets/a2a50cec-aac4-435e-b11c-262ab0e1d7b1" />

## Основные компоненты
Учебная система состоит из двух компонентов:
### Auth service — отвечает за “проверку доступа” (упрощённая логика).
-	выдаёт “токен” (упрощённо),
-	проверяет токен,
-	возвращает информацию: валиден/не валиден.

### Tasks service — CRUD задач, но каждый запрос требует проверки через Auth.
- хранит и управляет задачами,
-	перед выполнением операций проверяет токен через Auth.

## Коды статуса:
-	200 OK — успешный ответ
-	201 Created — ресурс создан
-	204 No Content — успешно, без тела
-	400 Bad Request — неверные данные
-	404 Not Found — ресурс не найден
-	422 Unprocessable Entity — некорректные данные по смыслу
-	500 Internal Server Error — внутренняя ошибка

# Примечания по конфигурации и требования

Для запуска требуется:

Go: версия 1.25.1

<img width="841" height="232" alt="Установка Git и Go" src="https://github.com/user-attachments/assets/8e01d831-5a7f-4376-8348-9052b240aec9" />


# Команды запуска/сборки
Для запуска http нужно выполнить 4 шага:
## 1) Клонировать данный репозиторий в удобную для вас папку:
```Powershell
git clone https://github.com/kayahan81/tech-ip-sem2
```
## 2) Перейти в папку http:
```Powershell
cd tech-ip-sem2
```
## 3) Загрузка зависимостей:
```Powershell
go mod tidy
```
## 4) Команда запуска
В первом окне
```Powershell
go run ./services/auth/cmd/auth
```
Во втором окне
```Powershell
go run ./services/tasks/cmd/tasks
```

# Проверка работоспособности
## Auth service
### POST /v1/auth/login
<img width="1632" height="86" alt="image" src="https://github.com/user-attachments/assets/a8449064-1176-4c0b-8f34-91c131d215eb" />

### GET /v1/auth/verify

## Tasks service
### POST /v1/tasks
<img width="1882" height="81" alt="image" src="https://github.com/user-attachments/assets/b68ada5e-4281-4bf6-9167-784f95cecba2" />

### GET /v1/tasks
<img width="1400" height="80" alt="image" src="https://github.com/user-attachments/assets/bb9ab81e-1a9d-48d4-92c3-563864c8cdb9" />

## GET /v1/tasks/{id}
<img width="1376" height="74" alt="image" src="https://github.com/user-attachments/assets/50d43d89-5c57-415b-b3c8-d93f988821e1" />

## PATCH /v1/tasks/{id}
Изменение и проверка.
<img width="1910" height="160" alt="image" src="https://github.com/user-attachments/assets/5b1039ca-af0d-4b83-84d6-604821df56b4" />

### Любая команда без токена
<img width="594" height="58" alt="image" src="https://github.com/user-attachments/assets/7b575da3-60cb-45ec-a3f9-67046ab4354d" />

# Ответы на вопросы
1.	Почему межсервисный вызов должен иметь таймаут?
Чтобы сервис не зависал навсегда в ожидании ответа от другого сервиса, если тот упал или недоступен по сети.
2.	Чем request-id помогает при диагностике ошибок?
Позволяет связать воедино логи разных сервисов, через которые прошел один запрос. По айди можно найти запросы и поэтапно пройти путь и найти где возникла ошибка.
3.	Какие статусы нужно вернуть клиенту при невалидном токене?
401 Unauthorized
4.	Чем опасно “делить одну БД” между сервисами?
Высокая нагрузка на базу данных
Изменение от одного сервиса может сломать взаимодействие с другим сервисом
