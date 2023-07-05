package redisDB

import (
	Interface "authorization/pkg/storage"
	"context"
	"github.com/redis/go-redis/v9"
	"log"
	"time"
)

// Storage Хранилище данных.
type Storage struct {
	db *redis.Client
}

// New Конструктор
func New(constr string) (*Storage, error) {
	client, err := redis.ParseURL(constr)
	if err != nil {
		log.Printf("неверная пользовательская информация %v\n", err)
	}
	s := &Storage{
		db: redis.NewClient(client),
	}
	return s, nil
}

// AddAccount Добавляет данные в базу redis.
func (s *Storage) AddAccount(c Interface.Account) error {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Set(ctx, c.Username, c.Password, 0).Err()
	if err != nil {
		log.Printf("Не удалось сделать запись %v\n", err)
		return nil
	}

	return nil
}

// SearchAccount Находит пароль по ключу в базе redis
func (s Storage) SearchAccount(c Interface.Account) (result string, err error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	result, err = s.db.Get(ctx, c.Username).Result()
	if err != nil {
		log.Printf("Нет такой записи %v\n", err)
		return "", err
	}
	return result, err
}

// KeysAccount Проверяет логин по ключу в базе redis
func (s Storage) KeysAccount(c Interface.Account) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	key, err := s.db.Keys(ctx, c.Username).Result()
	if err != nil {
		log.Printf("Ошибка при поиске ключей: %v\n", err)
		return false, err
	}

	// Обработка найденных ключей
	for _, k := range key {
		log.Println(k)
		if c.Username == k {
			return true, nil
		}
	}
	return false, nil
}

// DelAccount Удаляет аккаунт в базе Redis
func (s Storage) DelAccount(c Interface.Account) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()

	err := s.db.Del(ctx, c.Username).Err()
	if err != nil {
		log.Printf("Не удалось удвлить аккаунт %v\n", err)
		return false, nil
	}

	return true, nil
}
