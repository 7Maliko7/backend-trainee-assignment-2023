# Решение тестового задание для стажёра Backend

---

## Сервис динамического сегментирования пользователей

[Задача](https://github.com/avito-tech/backend-trainee-assignment-2023)

## Описание
- Язык программы Golang
- БД PostgreSQL

### Запуск
Запуск окружения с помощью `docker compose`. В директории [docker](/docker) расположен compose.yaml.

При первом запуске окружения БД пустая. Структура БД описана в файле [db.sql](/sql/db.sql).

#### Конфигурация
Для работы сервиса необходимо создать файл конфигурации. Пример в файле [config.yml.sample](/config.yml.sample).

Сервис сам прочитает конфигурацию из рабочей директории по умолчанию `./config.yml`. Иное расположение файла конфигурации можно указать опцией `-c` при запуске сервиса.

### API
API представляет собой REST.  
Коллекция [Postman](/postman/segments-api.postman_collection.json).  
Для запущенного сервиса Swagger расположен по адресу `{host}/docs/`.

---

#### POST /api/v1/segment/
Метод создания сегмента. Принимает slug (название) сегмента.

**Запрос**
``curl --location 'http://localhost:8080/api/v1/segment/' 
--header 'Content-Type: application/json' 
--data '{
"slug": "AVITO_DISCOUNT_30"
}'``  

**Ответ**
``{
"id": 7
}``

---

#### DELETE /api/v1/segment/{slug}
Метод удаления сегмента. Принимает slug (название) сегмента.

**Запрос**
``curl --location --request DELETE 'http://localhost:8080/api/v1/segment/AVITO_DISCOUNT_30'``

**Ответ**
``{
"status": "success"
}``

---

#### POST /api/v1/user/{user_id}/segments
Метод добавления/удаления пользователя в сегмент. Принимает список slug (названий) сегментов которые нужно добавить пользователю, список slug (названий) сегментов которые нужно удалить у пользователя, id пользователя.

Поле `segments` принимает объект типа `Segment`.

**Segment**

|Поле| Описание                                                  |
|----|-----------------------------------------------------------|
|slug| slug (название)                                           |
|action| Необходимое действие. Возможные значения: `add`, `delete` |

**Запрос**
``curl --location 'http://localhost:8080/api/v1/user/1001/segments' 
--header 'Content-Type: application/json' 
--data '{
"segments": [
{
"slug": "AVITO_DISCOUNT_30",
"action": "add"
},
{
"slug": "AVITO_PERFORMANCE_VAS",
"action": "delete"
}
]
}'``

**Ответ**
``{
"status": "success"
}``

---

#### GET /api/v1/user/{user_id}/segments
Метод получения активных сегментов пользователя. Принимает на вход id пользователя.

**Запрос**
``curl --location 'http://localhost:8080/api/v1/user/1000/segments'``

**Ответ**
``{
"segments": [
"AVITO_VOICE_MESSAGES"
]
}``

## Возникшие вопросы

- Многозначность метода добавления/удаления пользователя в сегмент.

Метод включает в себя действия добавления и удаления. Для выполнения однозначного действия был реализован объект `Segment`. Объект включает в себя указание действия, которое необходимо совершить над указанным сегментом для данного пользователя.