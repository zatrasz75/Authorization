package redisDB

import (
	"authorization/pkg/storage"
	"context"
	"testing"
	"time"
)

func TestNew(t *testing.T) {
	constr := "redis://localhost:6379"

	// Проверка ошибок
	s, err := New(constr)
	if err != nil {
		t.Fatalf("Ошибка при создании экземпляра Storage: %v", err)
	}

	// Проверка, что клиент Redis был успешно инициализирован
	if s.db == nil {
		t.Error("Клиент Redis не был инициализирован")
	}
}

func TestStorage_AddAccount(t *testing.T) {

	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	defer cancel()
	// Подготовка тестовых данных
	username := "testusers"
	password := "Test123!"
	account := storage.Account{
		Username: username,
		Password: password,
	}

	// Создание тестового экземпляра Storage
	storage, err := New("redis://localhost:6379")
	if err != nil {
		t.Fatalf("Ошибка при создании экземпляра Storage: %v", err)
	}

	// Вызов функции AddAccount
	err = storage.AddAccount(account)
	if err != nil {
		t.Fatalf("Ошибка при добавлении аккаунта: %v", err)
	}

	// Проверка, что аккаунт был успешно добавлен в базу данных
	// Проверка, что аккаунт был успешно добавлен в базу данных
	exists, err := storage.db.Exists(ctx, username).Result()
	if err != nil {
		t.Fatalf("Ошибка при проверке наличия ключа: %v", err)
	}
	if exists == 0 {
		t.Errorf("Аккаунт не найден в базе данных")
	}
}
