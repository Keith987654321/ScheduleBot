# Schedule bot
Schedule bot - бот для просмотра и управлением расписанием. 

## Требования
- Go (версия 1.18 или выше)
- PostgreSQL (версия 17.6 или совместимая)
- Docker (опционально, для запуска базы данных)
- Docker compose

## Установка и запуск
1. При помощи @BotFather в телеграмме создайте и скопируйте токен для бота.

2. В файл .env вставьте токен от своего бота в перменную окружения TOKEN.

3. Соберите и запустите докер-контейнер.
 - docker compose up --build  или  docker compose up --build -d, где -d (detached) чтобы контейнер работал в фоне.

4. Скопировать схему БД из файла db_scheme.sql
 - docker exec -i <container_id> psql -U postgres -d postgres < db_scheme.sql

5. Получить права администратора.
 - В телеграмме зайдите в чат с ботом и отправьте любое сообщение
 - docker exec -it <container_id> bash
 - psql -U postgres
 - UPDATE users SET role = 'admin';

6. С помощью терминала и sql-запросов или интерфейса бота в телеграмме добавьте расписание.