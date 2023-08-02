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
	"path/filepath"
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
	choice := os.Getenv("DEFINITION_DB")
	if choice == "" {
		log.Println("Не выбрана база данных ! Redis , Postgres или Mongo")
		return
	}

	// Можно сменить Порт при запуске флагом < --port-authorization= >
	portFlag := flag.String("port-authorization", port, "Порт для authorization сервиса")
	// Можно сменить Хост при запуске флагом < --host-authorization= >
	hostFlag := flag.String("host-authorization", host, "Хост для authorization сервиса")
	// Можно сменить URL соединения с бд при запуске флагом < --redis-url-authorization= >
	redis := flag.String("redis-url-authorization", dbRedis, "URL для соединения с Redis")
	// Можно сменить URL соединения с бд при запуске флагом < --postgres-url-authorization= >
	postgDB := flag.String("postgres-url-authorization", dbPostgres, "URL для соединения с Postgres")
	// Можно сменить URL соединения с бд при запуске флагом < --mongo-url-authorization= >
	mongo := flag.String("mongo-url-authorization", dbMongo, "URL для соединения с MongoDB")

	// Выбор базы данных Redis , Postgres или Mongo при запуске флагом < --select-db= >
	selectionDB := flag.String("select-db", choice, "Выбор базы данных Redis , Postgres или Mongo")

	flag.Parse()
	HOST := *hostFlag
	PORT := *portFlag
	REDIS := *redis
	POSGRES := *postgDB
	MONGO := *mongo
	CHOICE := *selectionDB

	// объект сервера
	var router server

	switch CHOICE {
	case "Redis":
		// объект базы данных Redis
		dbR, err := redisDB.New(REDIS)
		if err != nil {
			log.Printf("нет соединения с RedisDB %v", err)
			return
		}
		// Инициализируем хранилище сервера конкретной БД.
		router.db = dbR
	case "Postgres":
		// объект базы данных PostgreSQL
		dbP, err := postgres.New(POSGRES)
		if err != nil {
			log.Printf("нет соединения с PostgreSQL %v", err)
			return
		}
		// Инициализируем хранилище сервера конкретной БД.
		router.db = dbP

		//err = dbP.DropAccountsTable() // Удаляет таблицу с данными
		//if err != nil {
		//	log.Println(err)
		//}

		err = dbP.CreateAccountsTable() // Создает таблицу для данных
		if err != nil {
			log.Println(err)
		}
	case "Mongo":
		// объект базы данных Mongo
		dbM, err := mongoDB.New(MONGO)
		if err != nil {
			log.Printf("нет соединения с MongoDB %v", err)
			return
		}
		// Инициализируем хранилище сервера конкретной БД.
		router.db = dbM
	}

	// Получаем текущий путь к main.go
	currentDir, err := os.Getwd()
	if err != nil {
		log.Fatal("Не удалось получить текущий каталог:", err)
	}

	// Получаем абсолютный путь к каталогу web/
	webRoot := filepath.Join(currentDir, "../web")

	// Создаём объект API и регистрируем обработчики.
	router.api = api.New(router.db, webRoot)

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
		err := srv.ListenAndServe()
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
