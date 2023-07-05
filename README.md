# Authorization

### запустить приложение с переменными окружения по умолчанию:
* go run cmd/main.go

### изменить переменные через аргументы командной строки при запуске:
host
* go run cmd/main.go --host-authorization= < >

port
* go run cmd/main.go --port-authorization= < >

Redis URL
* go run cmd/main.go --rdis-url-authorization= < >

Postgres URL
* go run cmd/main.go --postgres-url-authorization= < >

### Или в файле .env
* APP_HOST , APP_PORT , DB_REDIS_URL , DB_POSTGRES_URL

### Доступные API для работы с базой данных Redis, примеры:

Регистрация, метод post (если успешно, перенаправляет авторизоваться)
* http://localhost:4000/registration

Авторизация, метод post перенаправляет на защищенную страницу
* http://localhost:4000/login

Защищенная страница, метод get (если не авторизован, перенаправляет авторизоваться)
* http://localhost:4000/dashboard

Удаление аккаунта, метод post
* http://localhost:4000/delaccount