# Final1lic
# Веб-сервис для вычисления арифметических выражений

## Описание
Этот проект реализует веб-сервис, который вычисляет арифметические выражения, переданные пользователем через HTTP-запрос.


## Запуск

1. Установите [Go](https://go.dev/dl/).
2. Установите [Git](https://git-scm.com/downloads).
3. Склонируйте проект с GitHub используя командную строку:
```
    git clone https://github.com/9doo/Final1lic
```
4. Перейдите в папку (Final1lic) проекта, выполните команду:
```
     go mod tidy
```
5. Запустите сервер:
```
    go run ./cmd/main.go
```
6. Сервис будет доступен по адресу: ```[http://localhost:8080/api/v1/calculate](http://localhost:8080/api/v1/calculate)```

### Как сменить порт (для Windiws)?
1. Для этого нужно собрать calc.exe 
2. Перейдите в папку (calc_go/cmd) проекта
3. выполните команды

```go build -o calc.exe``` 

```set "PORT=8087" & "calc.exe"``` 
В примере порт = 8087

## Эндпоинты
### `POST /api/v1/calculate`

#### Описание
Эндпоинт принимает JSON с математическим выражением.

#### Пример запроса с использованием curl
*Пример для cmd:*
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\\"expression\": \\"1\\"}"
```
*Пример корректного запроса, код:200*
```
git bash

curl --location 'localhost:8080/api/v1/calculate' \
--header 'Content-Type: application/json' \
--data '{
  "expression": "2+2*2"
}'
```
*Пример запроса с пустым выражением, код: 422, ошибка:empty expression*
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\\"expression\": \\"\\"}" 
```
*Пример запроса с делением на 0, код: 422, ошибка:division by zero*
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\\"expression\": \\"1/0\\"}" 
```
*Пример запроса с неверным выражением, код: 422, ошибка:invalid expression*
```
curl -X POST http://localhost:8080/api/v1/calculate -H "Content-Type: application/json" -d "{\\"expression\\": \\"1++*2\\"}" 
```
#### Для запросов можно использовать программу postman

## Команды для тестирования
Перейдите в каталог aplication или pkg\calculation и выполните команду:
```
  go test -v
```
