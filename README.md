# Тестовое задание golang-units

Перед запуском проекта требуется создать .env файл в корневом каталоге и заполнить его полями из примера в файле .env.example, ввести свой APOD_API_KEY

Доступ к методам http сервера можем получить по адресу: `/api/v1/apod_metadata`

- если указать параметр date в формате "2023-08-22", то получим запись за определенную дату
- если не указывать параметры, то получим все записи

Для деплоя проекта использовать команду: `make deploy`
