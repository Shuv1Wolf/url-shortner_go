# url-shortner

Проект REST API сделанный в целях обучения и демонстрации возможностей.

## v1.0.0

### Эндпоинты

- /url (POST, DELETE) - добавление и удаление ссылок, принимает тело запроса в формате JSON и возвращает статус 200. Требуется Basic Auth авторизация.

Пример тела запроса:
```
(POST, добавление ссылки)
{
    "url": "https://translate.google.com/",
    "alias": "translate"
}
```
```
(DELETE, удаление ссылки)
{
    "alias": "translate"
}
```
Примеры ответа:
```
{
    "status": "OK",
    "alias": "translate"
}
```
```
{
    "status": "ERROR",
    "error": "url already exists"
}
``` 
```
{
    "status": "ERROR",
    "error": "failed decode request"
}
```
и тд...

- /{alias} (GET) - принимает алиас ссылки и возвращает статус 302 или 200. В случае успеха перенаправляет на страницу.

Пример ответа с ошибкой:
```
{
    "status": "ERROR",
    "error": "not found"
}
```

### Плюшки

- есть кастомный логгер для вывода сообщений в консоль при env: "local":

<p>
 <img width="500px" src="img/logger.png" alt="logger"/>
</p>

### На будущее

- добавить ORM
- добавить JWT авторизацию





