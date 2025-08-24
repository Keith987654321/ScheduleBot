# Schedule bot
Schedule bot - бот для просмотра и управлением расписанием. 

## Требования
- Go (версия 1.18 или выше)
- PostgreSQL (версия 17.6 или совместимая)
- Docker (опционально, для запуска базы данных)
 
## Установка и запуск
1. Создайте бд postgresSQL по схеме из файла scheme.sql
    psql -U <user> -d <database> -f schema.sql

    или с помощью Docker: 
    docker run --name postgres -e POSTGRES_USER=<user> -e POSTGRES_PASSWORD=<pass> -e POSTGRES_DB=<name> -p 5432:5432 -d postgres:17.6

2. Установите все зависимости с помощью команды в терминале go mod tidy

3. При помощи @BotFather в телеграмме создайте и скопируйте токен для бота

4. Запустите бота с помощью команды go run main.go -user=<username> -pass=<password> -name=<database_name> -ssl=<enable/disable> -token=<telegram_token>

    - Флаг user - имя пользователя бд
    - Флаг pass - пароль бд
    - Флаг name - имя бд
    - Флаг ssl - режим ssl (enable / disable)
    - Флаг token - токен вашего бота в тг, который вы получили на 3 шаге

    Пример Команды:
    go run main.go -user=postgres -pass=secret -name=schedule_db -ssl=disable -token=123456:ABC-DEF1234ghIkl-zyx57W2v1u123ew11