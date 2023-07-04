package api_test

import (
	"authorization/pkg/api"
	"authorization/pkg/storage"
	"authorization/pkg/storage/redisDB"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

func TestRegistrationHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаем тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаем экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаем данные регистрации
	formData := url.Values{}
	formData.Set("username", "ups@mail.ru")
	formData.Set("password", "Test123!")

	// Создаем запрос POST с данными регистрации
	req, err := http.NewRequest(http.MethodPost, "/registration", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Создаем ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusFound)
	}

	// Проверяем ожидаемое сообщение
	expectedMessage := "Ваш аккаунт успешно создан."
	actualMessage := resRecorder.Body.String()
	if actualMessage != expectedMessage {
		t.Errorf("Неверное сообщение: получено %v, ожидается %v", actualMessage, expectedMessage)
	}

	// Проверяем, что аккаунт был добавлен в базу данных
	account := storage.Account{
		Username: "ups@mail.ru",
		Password: "Test123!",
	}
	keys, err := db.KeysAccount(account)
	if err != nil {
		t.Fatalf("Ошибка при поиске ключей: %v", err)
	}
	if !keys {
		t.Errorf("Аккаунт не найден в базе данных")
	}
}

func TestAPI_loginHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаем тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаем экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаем данные регистрации
	formData := url.Values{}
	formData.Set("username", "ups@mail.ru")
	formData.Set("password", "Test123!")

	// Создаем запрос POST с данными регистрации
	req, err := http.NewRequest(http.MethodPost, "/login", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Создаем ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusFound {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusFound)
	}

}

func TestDashboardHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаем тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаем экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаем запрос GET для защищенной страницы
	req, err := http.NewRequest(http.MethodGet, "/dashboard", nil)
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}

	// Устанавливаем cookie авторизованного пользователя
	cookie := &http.Cookie{
		Name:  "session",
		Value: "authenticated",
	}
	req.AddCookie(cookie)

	// Создаем ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusOK)
	}

	// Проверяем ожидаемое сообщение
	expectedMessage := "Добро пожаловать в панель управления!"
	actualMessage := resRecorder.Body.String()
	if actualMessage != expectedMessage {
		t.Errorf("Неверное сообщение: получено %v, ожидается %v", actualMessage, expectedMessage)
	}
}

func TestAPI_delAccountHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаем тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаем экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаем данные регистрации
	formData := url.Values{}
	formData.Set("username", "ups@mail.ru")

	// Создаем запрос POST с данными регистрации
	req, err := http.NewRequest(http.MethodPost, "/delaccount", strings.NewReader(formData.Encode()))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// Создаем ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusOK)
	}

	// Проверяем ожидаемое сообщение
	expectedMessage := "Ваш аккаунт успешно удален."
	actualMessage := resRecorder.Body.String()
	if actualMessage != expectedMessage {
		t.Errorf("Неверное сообщение: получено %v, ожидается %v", actualMessage, expectedMessage)
	}

}
