package main

import (
	"authorization/pkg/api"
	"authorization/pkg/middl"
	"authorization/pkg/storage"
	"authorization/pkg/storage/redisDB"
	"log"
	"net/http"
)

// сервер.
type server struct {
	db  storage.Interface
	api *api.API
}

func main() {

	// объект сервера
	var srv server

	constr := "redis://localhost:6379"

	// объект базы данных redis
	db, err := redisDB.New(constr)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	srv.api.Router().Use(middl.Middle)

	log.Println("Запуск сервера на http://127.0.0.1:4000/login")

	err = http.ListenAndServe(":4000", srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
