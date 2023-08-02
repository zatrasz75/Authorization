package api

import (
	"authorization/pkg/check"
	"authorization/pkg/storage"
	"encoding/json"
	"github.com/gorilla/mux"
	"io/ioutil"
	"log"
	"net/http"
	"path/filepath"
	"strings"
	"sync"
)

// API приложения.
type API struct {
	r       *mux.Router       // Маршрутизатор запросов
	db      storage.Interface // база данных
	webRoot string            // Корневая директория для веб-приложения
}

// New Конструктор API.
func New(db storage.Interface, webRoot string) *API {
	api := API{
		r:       mux.NewRouter(),
		db:      db,
		webRoot: webRoot,
	}
	//	api.r = mux.NewRouter()
	api.endpoints()
	return &api
}

// Router возвращает маршрутизатор запросов.
func (api *API) Router() *mux.Router {
	return api.r
}

// Регистрация обработчиков API.
func (api *API) endpoints() {
	api.r.HandleFunc("/", api.home).Methods(http.MethodGet)
	api.r.HandleFunc("/api/login", api.handleLogin).Methods(http.MethodGet)
	api.r.HandleFunc("/login", api.loginHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/dashboard", api.dashboardHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/registration", api.registrationHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/delaccount", api.delAccountHandler).Methods(http.MethodPost)

	// веб-приложение
	api.r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./web/"))))

}

// Обработчик для статических файлов веб-приложения.
func (api *API) serveWebFiles(w http.ResponseWriter, r *http.Request) {
	filePath := r.URL.Path
	// Проверяем, что запрошенный путь начинается с "/web/".
	if !strings.HasPrefix(filePath, "/web/") {
		http.NotFound(w, r)
		return
	}

	// Проверяем, что путь после "/web/" не содержит "../" (попытка обхода пути).
	if strings.Contains(filePath, "../") {
		http.NotFound(w, r)
		return
	}

	// Строим абсолютный путь к файлу.
	absolutePath := filepath.Join(api.webRoot, filePath[5:])

	// Проверяем, что абсолютный путь находится в пределах корневой директории для веб-приложения.
	if !strings.HasPrefix(absolutePath, api.webRoot) {
		http.NotFound(w, r)
		return
	}

	// Обслуживаем статический файл.
	http.ServeFile(w, r, absolutePath)
}

func (api *API) handleLogin(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/api/login" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Чтение содержимого файла login.html
	content, err := ioutil.ReadFile("web/registration.html")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// Отправка содержимого файла как ответ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

func (api *API) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	// Чтение содержимого файла login.html
	content, err := ioutil.ReadFile("web/login.html")
	if err != nil {
		http.Error(w, "Ошибка чтения файла", http.StatusInternalServerError)
		return
	}

	// Отправка содержимого файла как ответ
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write(content)
}

// Функция-обработчик для страницы с регистрацией
func (api *API) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/registration" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из формы регистрации
	var f storage.FormAccount
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	// Каналы для синхронизации и передачи результатов проверок
	emailResultCh := make(chan bool, 1)
	letterCh := make(chan bool, 1)
	specCharCh := make(chan bool, 1)
	lenRegexCh := make(chan bool, 1)
	numbersCh := make(chan bool, 1)
	containLetterCh := make(chan bool, 1)
	weakCh := make(chan bool, 1)

	var wg sync.WaitGroup
	wg.Add(7) // Устанавливаем количество ожидаемых горутин

	// Горутина для проверки адреса электронной почты
	go func() {
		defer wg.Done()
		emailResultCh <- check.CheckEmail(f.Username)
	}()
	// Горутины для проверки пароля
	go func() {
		defer wg.Done()
		letterCh <- check.LowercaseLetter(f.Password)
	}()
	go func() {
		defer wg.Done()
		specCharCh <- check.SpecCharRegex(f.Password)
	}()
	go func() {
		defer wg.Done()
		lenRegexCh <- check.LenPass(f.Password)
	}()
	go func() {
		defer wg.Done()
		numbersCh <- check.NumbersPass(f.Password)
	}()
	go func() {
		defer wg.Done()
		containLetterCh <- check.ContainPass(f.Password)
	}()
	go func() {
		defer wg.Done()
		weakCh <- check.WeakPass(f.Password)
	}()
	wg.Wait()

	var errorMessages []string

	if !<-emailResultCh {
		errorMessages = append(errorMessages, "Адрес электронной почты не корректный")
	}
	if !<-letterCh {
		errorMessages = append(errorMessages, "Ошибка! Пароль должен содержать строчные буквы")
	}
	if !<-specCharCh {
		errorMessages = append(errorMessages, "Ошибка! Пароль должен содержать спец. символ")
	}
	if !<-lenRegexCh {
		errorMessages = append(errorMessages, "Ошибка! Пароль должен содержать не менее 8 символов")
	}
	if !<-numbersCh {
		errorMessages = append(errorMessages, "Ошибка! Пароль должен содержать цифры")
	}
	if !<-containLetterCh {
		errorMessages = append(errorMessages, "Ошибка! Пароль должен содержать прописные буквы")
	}
	if !<-weakCh {
		errorMessages = append(errorMessages, "Предупреждение! Очень слабый пароль, придумайте другой")
	}

	// Вывод ошибок, если они есть
	if len(errorMessages) > 0 {
		resp := storage.Response{
			Success:       false,
			ErrorMessages: errorMessages,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		return
	}

	// Адрес электронной почты и пароль валидны
	hash := check.HashPass(f.Password)
	c := storage.Account{
		Username: f.Username,
		Password: hash,
	}

	// Проверяем есть ли такой пользователь в базе данных
	keys, err := api.db.KeysAccount(c)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка при проверке пользователя", http.StatusInternalServerError)
		return
	}

	if keys == true {
		resp := storage.Response{
			Success: false,
			Message: "Такой пользователь уже существует",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)

		return
	} else {
		err = api.db.AddAccount(c)
		if err != nil {
			log.Println(err)
			http.Error(w, "Ошибка при добавлении пользователя", http.StatusInternalServerError)
			return
		}
		resp := storage.Response{
			Success: true,
			Message: "Ваш аккаунт успешно создан.",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
		return
	}
}

// Функция-обработчик для страницы с авторизацией
func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
	}
	// Получаем данные из формы авторизации
	var f storage.FormAccount
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}

	hash := check.HashPass(f.Password)
	c := storage.Account{
		Username: f.Username,
		Password: hash,
	}

	// Получаем пароль из базы данных
	result, err := api.db.SearchAccount(c)
	if err != nil {
		log.Println(err)
	}

	// Проверяем, соответствуют ли переданные данные ожидаемым значениям
	if c.Password == result {
		// Если авторизация успешна, сохраняем информацию о входе в сессии
		session, err := r.Cookie("session")
		if err != nil {
			session = &http.Cookie{
				Name:   "session",
				Value:  "authenticated",
				Secure: true,
				MaxAge: 3600, // 3600 Устанавливает время жизни cookie на 1 час

			}
		}
		http.SetCookie(w, session)

		// Перенаправляем пользователя на защищенную страницу
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		// Если авторизация не удалась, отображаем сообщение об ошибке
		resp := storage.Response{
			Success: false,
			Message: "Нет такой записи, проверти логин или пароль",
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

// Функция-обработчик для защищенной страницы
func (api *API) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
		return
	}
	// Проверяем, авторизован ли пользователь
	session, err := r.Cookie("session")
	if err != nil || session.Value != "authenticated" {
		// Если пользователь не авторизован, перенаправляем на страницу входа
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Если пользователь авторизован, отображаем JSON-ответ и содержимое файла dashboard.html
	resp := storage.Response{
		Success: true,
		Message: "Добро пожаловать в панель управления !!!",
	}

	//tmpl, err := template.ParseFiles("web/ups.html")
	//if err != nil {
	//	log.Println("templateОшибка при обработке шаблона:", err)
	//	http.Error(w, "Ошибка при обработке шаблона", http.StatusInternalServerError)
	//	return
	//}
	//// Устанавливаем правильный Content-Type для HTML
	//w.Header().Set("Content-Type", "text/html; charset=utf-8")
	//
	//err = tmpl.Execute(w, resp)
	//if err != nil {
	//	log.Println("Execute Ошибка при выполнении шаблона:", err)
	//	http.Error(w, "Ошибка при выполнении шаблона", http.StatusInternalServerError)
	//	return
	//}

	// Отправляем JSON-ответ
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}

// Функция-обработчик для страницы с удалением аккаунта
func (api *API) delAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delaccount" {
		http.NotFound(w, r)
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Метод не поддерживается", http.StatusMethodNotAllowed)
		return
	}

	// Получаем данные из формы удаления
	var f storage.FormAccount
	if err := json.NewDecoder(r.Body).Decode(&f); err != nil {
		http.Error(w, "Ошибка при декодировании JSON", http.StatusBadRequest)
		return
	}
	c := storage.Account{
		Username: f.Username,
	}

	// Проверяем есть ли такой пользователь в базе redis
	keys, err := api.db.KeysAccount(c)
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка при проверке пользователя", http.StatusInternalServerError)
		return
	}

	// Удаляем аккаунт, если он существует
	if keys == true {
		a, err := api.db.DelAccount(c)
		if err != nil {
			log.Println(err)
			http.Error(w, "Ошибка при удалении пользователя", http.StatusInternalServerError)
			return
		}
		if a == true {
			// Удаляем Cookie
			sessionCookie := &http.Cookie{
				Name:   "session",
				Value:  "",
				MaxAge: -1, // или 0
				Path:   "/",
			}
			http.SetCookie(w, sessionCookie)

			resp := storage.Response{
				Success: true,
				Message: "Ваш аккаунт успешно удален.",
			}
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(resp)
			return
		}
	}
	// Если аккаунт не существует
	resp := storage.Response{
		Success: false,
		Message: "Такой пользователь не существует, проверьте логин.",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(resp)
}
