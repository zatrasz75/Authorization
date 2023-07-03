package main

import (
	"authorization/pkg/api"
	"authorization/pkg/middl"
	"authorization/pkg/storage"
	"authorization/pkg/storage/redisDB"
	"flag"
	"log"
	"net/http"
)

// сервер.
type server struct {
	db  storage.Interface
	api *api.API
}

const (
	authorizationPort = ":4000"
)

func main() {

	// Можно сменить Порт при запуске флагом < --port-authorization= >
	portFlag := flag.String("port-authorization", authorizationPort, "Порт для authorization сервиса")
	flag.Parse()
	PORT := *portFlag

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

	log.Println("Запуск сервера на http://127.0.0.1" + PORT + "/login")

	err = http.ListenAndServe(PORT, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
