package main

import (
	"authorization/pkg/api"
	"authorization/pkg/middl"
	"authorization/pkg/storage"
	"authorization/pkg/storage/mongoDB"
	"authorization/pkg/storage/postgres"
	"authorization/pkg/storage/redisDB"
	"context"
	"flag"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// сервер
type server struct {
	db  storage.Interface
	api *api.API
}

// init вызывается перед main()
func init() {
	// загружает значения из файла .env в систему
	if err := godotenv.Load(); err != nil {
		log.Print("Файл .env не найден.")
	}
}

// Переменные для окружения соединения
const (
	authorizationPort = "5000"
	authorizationHost = "127.0.0.1"
	clientRedisDB     = "redis://localhost:6379"
	clientPostgresDB  = "postgres://postgres:postgrespw@localhost:49153/Account"
	clientMongoDB     = "mongodb://localhost:27015/"
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
	dbPostgres := os.Getenv("DB_POSTGRES_URL")
	if dbPostgres == "" {
		dbPostgres = clientPostgresDB
	}
	dbMongo := os.Getenv("DB_MONGO_URL")
	if dbMongo == "" {
		dbMongo = clientMongoDB
	}

	// Можно сменить Порт при запуске флагом < --port-authorization= >
	portFlag := flag.String("port-authorization", port, "Порт для authorization сервиса")
	// Можно сменить Хост при запуске флагом < --host-authorization= >
	hostFlag := flag.String("host-authorization", host, "Хост для authorization сервиса")
	// Можно сменить URL соединения с бд при запуске флагом < --rdis-url-authorization= >
	redis := flag.String("rdis-url-authorization", dbRedis, "URL для соединения с Redis")
	// Можно сменить URL соединения с бд при запуске флагом < --postgres-url-authorization= >
	postgDB := flag.String("postgres-url-authorization", dbPostgres, "URL для соединения с Postgres")
	// Можно сменить URL соединения с бд при запуске флагом < --mongo-url-authorization= >
	mongo := flag.String("mongo-url-authorization", dbMongo, "URL для соединения с MongoDB")

	flag.Parse()
	HOST := *hostFlag
	PORT := *portFlag
	REDIS := *redis
	POSGRES := *postgDB
	MONGO := *mongo

	// объект базы данных Redis
	dbR, err := redisDB.New(REDIS)
	if err != nil {
		log.Printf("нет соединения с RedisDB %v", err)
	}

	// объект базы данных PostgreSQL
	dbP, err := postgres.New(POSGRES)
	if err != nil {
		log.Printf("нет соединения с PostgreSQL %v", err)
	}

	// объект базы данных Mongo
	dbM, err := mongoDB.New(MONGO)
	if err != nil {
		log.Printf("нет соединения с MongoDB %v", err)
	}

	_, _, _ = dbR, dbP, dbM

	// объект сервера
	var router server

	// Инициализируем хранилище сервера конкретной БД.
	router.db = dbP

	if router.db == dbP {
		//err = dbP.DropAccountsTable()
		//if err != nil {
		//	log.Println(err)
		//}

		err = dbP.CreateAccountsTable()
		if err != nil {
			log.Println(err)
		}
	}

	// Создаём объект API и регистрируем обработчики.
	router.api = api.New(router.db)

	router.api.Router().Use(middl.Middle)

	log.Println("Запуск сервера на ", "http://"+HOST+":"+PORT)

	// Создаем HTTP сервер с заданным адресом и обработчиком.
	srv := http.Server{
		Addr:         ":" + PORT,
		Handler:      router.api.Router(),
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 30 * time.Second,
	}

	go func() {
		err = srv.ListenAndServe()
		if err != nil {
			log.Fatal("Не удалось запустить сервер шлюза. Error:", err)
		}
	}()

	graceShutdown(srv)
}

// Выключает сервер
func graceShutdown(srv http.Server) {
	quitCH := make(chan os.Signal, 1)
	signal.Notify(quitCH, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	<-quitCH

	// Создаем контекст с таймаутом 5 секунд.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	// Останавливаем сервер с таймаутом 5 секунд.
	err := srv.Shutdown(ctx)
	if err != nil {
		log.Printf("Ошибка при закрытии прослушивателей или тайм-аут контекста %v", err)
		return
	}
	log.Printf("Выключение сервера")
	os.Exit(0)
}
