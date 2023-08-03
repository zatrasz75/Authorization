# Authorization

### запустить приложение с переменными окружения по умолчанию и выбором базы данных Redis , Postgres или Mongo:
#### Пример с Postgres
* go run cmd/main.go --select-db=Postgres

### изменить переменные через аргументы командной строки при запуске:
host
* go run cmd/main.go --host-authorization= < >

port
* go run cmd/main.go --port-authorization= < >

Redis URL
* go run cmd/main.go --rdis-url-authorization= < >

Postgres URL
* go run cmd/main.go --postgres-url-authorization= < >

Mongo URL
* go run cmd/main.go --mongo-url-authorization= < >

### Или в файле .env
* APP_HOST , APP_PORT , DB_REDIS_URL , DB_POSTGRES_URL , DB_MONGO_URL , DEFINITION_DB

### Доступные API для работы с выбранной базой данных , примеры:

Регистрация, метод post (если успешно, перенаправляет авторизоваться)
* http://127.0.0.1:5000/api/login

Авторизация, метод post перенаправляет на защищенную страницу
* http://localhost:5000/

#### Для Postman
Защищенная страница, метод get (если не авторизован, возвращает ошибку)
* http://localhost:5000/dashboard/

Удаление аккаунта, метод post
* http://localhost:5000/delaccount