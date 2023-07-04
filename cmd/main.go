package main

import (
	"authorization/pkg/api"
	"authorization/pkg/middl"
	"authorization/pkg/storage"
	"authorization/pkg/storage/redisDB"
	"flag"
	"log"
	"net/http"
	"os"
)

// сервер
type server struct {
	db  storage.Interface
	api *api.API
}

// Переменные для окружения соединения
const (
	authorizationPort = ":4000"
	authorizationHost = "127.0.0.1"
	clientRedisDB     = "redis://localhost:6379"
)

func main() {
	port := os.Getenv("APP_PORT")
	if port == "" {
		port = authorizationPort
	}
	host := os.Getenv("APP_HOST")
	if host == "" {
		host = authorizationHost
	}
	dbRedis := os.Getenv("DB_REDIS_URL")
	if dbRedis == "" {
		dbRedis = clientRedisDB
	}

	// Можно сменить Порт при запуске флагом < --port-authorization= >
	portFlag := flag.String("port-authorization", port, "Порт для authorization сервиса")
	// Можно сменить Хост при запуске флагом < --host-authorization= >
	hostFlag := flag.String("host-authorization", host, "Хост для authorization сервиса")
	// Можно сменить URL соединения с бд при запуске флагом < --rdis-url-authorization= >
	redis := flag.String("rdis-url-authorization", dbRedis, "URL для соединения с ")
	flag.Parse()
	HOST := *hostFlag
	PORT := *portFlag
	REDIS := *redis

	// объект сервера
	var srv server

	// объект базы данных redis
	db, err := redisDB.New(REDIS)
	if err != nil {
		log.Fatal(err)
	}
	// Инициализируем хранилище сервера конкретной БД.
	srv.db = db

	// Создаём объект API и регистрируем обработчики.
	srv.api = api.New(srv.db)

	srv.api.Router().Use(middl.Middle)

	log.Println("Запуск сервера на ", "http://"+HOST+PORT+"/login")

	err = http.ListenAndServe(HOST+PORT, srv.api.Router())
	if err != nil {
		log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
	}

}
