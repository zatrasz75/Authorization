package api_test

import (
	"authorization/pkg/api"
	"authorization/pkg/storage"
	"authorization/pkg/storage/redisDB"
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestRegistrationHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаём тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаём экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаём данные регистрации
	formData := storage.FormAccount{
		Username: "ups@mail.ru",
		Password: "Test123!",
	}
	jsonData, err := json.Marshal(formData)
	if err != nil {
		t.Fatalf("Ошибка при преобразовании данных в JSON: %v", err)
	}

	// Создаём запрос POST с данными регистрации
	req, err := http.NewRequest(http.MethodPost, "/registration", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаём ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusFound)
	}

	// Проверяем ожидаемый JSON-ответ
	expectedJSON := `{"success":true,"message":"Ваш аккаунт успешно создан.","errorMessages":null}`
	actualJSON := strings.TrimSpace(resRecorder.Body.String())
	if actualJSON != expectedJSON {
		t.Errorf("Неверный JSON-ответ: получено\n%s\nожидается\n%s", actualJSON, expectedJSON)
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
	// Создаём тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаём экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаём данные регистрации
	formData := storage.FormAccount{
		Username: "ups@mail.ru",
		Password: "Test123!",
	}
	jsonData, err := json.Marshal(formData)
	if err != nil {
		t.Fatalf("Ошибка при преобразовании данных в JSON: %v", err)
	}

	// Создаём запрос POST с данными регистрации
	req, err := http.NewRequest(http.MethodPost, "/login", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаём ResponseWriter для записи ответа
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
	// Создаём тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаём экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаём запрос GET для защищённой страницы
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

	// Создаём ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusOK)
	}

	// Проверяем ожидаемое сообщение
	expectedMessage := `{"success":true,"message":"Добро пожаловать в панель управления !!!","errorMessages":null}`
	actualMessage := strings.TrimSpace(resRecorder.Body.String())
	if actualMessage != expectedMessage {
		t.Errorf("Неверное сообщение: получено %v, ожидается %v", actualMessage, expectedMessage)
	}
}

func TestAPI_delAccountHandler(t *testing.T) {
	constr := "redis://localhost:6379"
	// Создаём тестовую базу данных
	db, _ := redisDB.New(constr)

	// Создаём экземпляр API с тестовой базой данных
	a := api.New(db)

	// Создаём аккаунт для удаления
	formData := storage.FormAccount{
		Username: "ups@mail.ru",
	}
	jsonData, err := json.Marshal(formData)
	if err != nil {
		t.Fatalf("Ошибка при преобразовании данных в JSON: %v", err)
	}

	// Создаём запрос POST с данными удаления
	req, err := http.NewRequest(http.MethodPost, "/delaccount", bytes.NewBuffer(jsonData))
	if err != nil {
		t.Fatalf("Ошибка при создании запроса: %v", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Создаём ResponseWriter для записи ответа
	resRecorder := httptest.NewRecorder()

	// Выполняем запрос
	a.Router().ServeHTTP(resRecorder, req)

	// Проверяем статус код
	if resRecorder.Code != http.StatusOK {
		t.Errorf("Неверный статус код: получено %v, ожидается %v", resRecorder.Code, http.StatusOK)
	}

	// Проверяем ожидаемый JSON-ответ при успешном удалении
	expectedResponse := `{"success":true,"message":"Ваш аккаунт успешно удален.","errorMessages":null}`
	actualResponse := strings.TrimSpace(resRecorder.Body.String())
	if actualResponse != expectedResponse {
		t.Errorf("Неверный JSON-ответ: получено %v, ожидается %v", actualResponse, expectedResponse)
	}

	// Создаём аккаунт для проверки удаления из базы данных
	c := storage.Account{
		Username: "testuser",
	}

	// Проверяем, что аккаунт был удалён из базы данных
	keys, err := db.KeysAccount(c)
	if err != nil {
		t.Fatalf("Ошибка при поиске ключей: %v", err)
	}
	if keys {
		t.Errorf("Аккаунт не должен быть найден в базе данных после удаления")
	}
}
