# url-shortner

Проект REST API сделанный в целях обучения и демонстрации возможностей.

### v1.0.0

- /url (POST, DELETE) - добавление и удаление ссылок, принимает тело запроса в формате JSON и возвращает статус 200 или 400 в случае ошибки. Требуется Basic Auth авторизация.

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



