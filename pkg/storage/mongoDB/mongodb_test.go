package mongoDB

import (
	"authorization/pkg/storage"
	"testing"
)

func TestNew(t *testing.T) {
	s, err := New("mongodb://localhost:27015/")
	if err != nil {
		t.Fatal(err)
	}

	// Проверка, что клиент MongoDB был успешно инициализирован
	if s.db == nil {
		t.Error("Клиент MongoDB не был инициализирован")
	}
}

func TestStorage_AddAccount(t *testing.T) {
	// Установка соединения с базой данных MongoDB
	dataBase, err := New("mongodb://localhost:27015/")
	if err != nil {
		t.Fatalf("не удалось подключиться к базе данных: %v", err)
	}
	// Определение тестового аккаунта
	c := storage.Account{
		Username: "krex@ya.ru",
		Password: "12345678",
	}
	// Вызов функции AddAccount
	err = dataBase.AddAccount(c)
	if err != nil {
		t.Error(err)
	}

	t.Log("Запись создана.")
}

func TestStorage_SearchAccount(t *testing.T) {
	// Установка соединения с базой данных MongoDB
	dataBase, err := New("mongodb://localhost:27015/")
	if err != nil {
		t.Fatalf("не удалось подключиться к базе данных: %v", err)
	}

	// Определение тестового аккаунта
	c := storage.Account{
		Username: "krex@ya.ru",
	}

	// Вызов функции SearchAccount
	password, err := dataBase.SearchAccount(c)
	if err != nil {
		t.Errorf("ошибка при поиске аккаунта: %v", err)
		return
	}

	// Проверка полученного пароля
	expectedPassword := "12345678" // Замените на ожидаемый пароль
	if password != expectedPassword {
		t.Errorf("неправильный пароль. Получено: %s, Ожидается: %s", password, expectedPassword)
	}
}

func TestStorage_KeysAccount(t *testing.T) {
	// Установка соединения с базой данных MongoDB
	dataBase, err := New("mongodb://localhost:27015/")
	if err != nil {
		t.Fatalf("не удалось подключиться к базе данных: %v", err)
	}

	// Определение тестового аккаунта
	c := storage.Account{
		Username: "krex@ya.ru",
	}

	// Вызов функции KeysAccount
	exists, err := dataBase.KeysAccount(c)
	if err != nil {
		t.Errorf("ошибка при проверке наличия аккаунта: %v", err)
		return
	}

	// Проверка результата наличия аккаунта
	if exists != true {
		t.Errorf("неправильный результат проверки наличия аккаунта. Получено: %v, Ожидается: %v", exists, true)
	}
}

func TestStorage_DelAccount(t *testing.T) {
	// Установка соединения с базой данных MongoDB
	dataBase, err := New("mongodb://localhost:27015/")
	if err != nil {
		t.Fatalf("не удалось подключиться к базе данных: %v", err)
	}

	// Определение тестового аккаунта
	c := storage.Account{
		Username: "krex@ya.ru",
	}

	// Вызов функции DelAccount
	deleted, err := dataBase.DelAccount(c)
	if err != nil {
		t.Errorf("ошибка при удалении аккаунта: %v", err)
		return
	}

	// Проверка результата удаления аккаунта
	//expectedDeleted := true // Замените на ожидаемое значение
	if deleted != true {
		t.Errorf("неправильный результат удаления аккаунта. Получено: %v, Ожидается: %v", deleted, true)
	}

}
