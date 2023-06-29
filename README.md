# Authorization

### запустить приложение:
* go run cmd/main.go

### Доступные API для работы с базой данных Redis, примеры:

Регистрация, метод post (если успешно, перенаправляет авторизоваться)
* http://localhost:4000/registration

Авторизация, метод post перенаправляет на защищенную страницу
* http://localhost:4000/login

Защищенная страница, метод get (если не авторизован, перенаправляет авторизоваться)
* http://localhost:4000/dashboard

Удаление аккаунта, метод post
* http://localhost:4000/delaccount