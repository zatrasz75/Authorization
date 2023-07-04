package api

import (
	"authorization/pkg/check"
	"authorization/pkg/storage"
	"fmt"
	"github.com/gorilla/mux"
	"html/template"
	"log"
	"net/http"
)

// API приложения.
type API struct {
	r  *mux.Router       // Маршрутизатор запросов
	db storage.Interface // база данных
}

// New Конструктор API.
func New(db storage.Interface) *API {
	api := API{
		r:  mux.NewRouter(),
		db: db,
	}
	api.r = mux.NewRouter()
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
	api.r.HandleFunc("/login", api.loginHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/dashboard", api.dashboardHandler).Methods(http.MethodGet)
	api.r.HandleFunc("/registration", api.registrationHandler).Methods(http.MethodPost)
	api.r.HandleFunc("/delaccount", api.delAccountHandler).Methods(http.MethodPost)
	// веб-приложение
	api.r.PathPrefix("/web/").Handler(http.StripPrefix("/web/", http.FileServer(http.Dir("./cmd/web"))))

}

func (api *API) home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	// Определение данных для передачи в шаблон
	data := struct {
		Title string
		Text  string
	}{
		Title: "Форма авторизации",
		Text:  "Ups",
	}

	// Загрузка и компиляция шаблона
	tmpl, err := template.ParseFiles("./cmd/web/index.html")
	if err != nil {
		log.Println(err)
		http.Error(w, "Ошибка открытия файла", http.StatusInternalServerError)
		return
	}

	// Рендеринг шаблона с передачей данных
	err = tmpl.Execute(w, data)
	if err != nil {
		http.Error(w, "Ошибка сервера", http.StatusInternalServerError)
		return
	}

}

// Функция-обработчик для страницы с регистрацией
func (api *API) registrationHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/registration" {
		http.NotFound(w, r)
	}
	// Получаем данные из формы регистрации
	username := r.FormValue("username")
	password := r.FormValue("password")

	// Каналы для синхронизации и передачи результатов проверок
	emailResultCh := make(chan bool)
	letterCh := make(chan bool)
	specCharCh := make(chan bool)
	lenRegexCh := make(chan bool)
	numbersCh := make(chan bool)
	containLetterCh := make(chan bool)
	weakCh := make(chan bool)

	// Горутина для проверки адреса электронной почты
	go func() {
		emailResultCh <- check.CheckEmail(username)
	}()
	// Горутина для проверки пароля
	go func() {
		letterCh <- check.LowercaseLetter(password)
	}()
	go func() {
		specCharCh <- check.SpecCharRegex(password)
	}()
	go func() {
		lenRegexCh <- check.LenPass(password)
	}()
	go func() {
		numbersCh <- check.NumbersPass(password)
	}()
	go func() {
		containLetterCh <- check.ContainPass(password)
	}()
	go func() {
		weakCh <- check.WeakPass(password)
	}()

	// Ожидание результатов проверок
	emailValid := <-emailResultCh
	letterValid := <-letterCh
	specValid := <-specCharCh
	lenValid := <-lenRegexCh
	numbersValid := <-numbersCh
	containValid := <-containLetterCh
	weakValid := <-weakCh

	// Проверка результатов
	if !emailValid {
		fmt.Fprintf(w, "Адрес электронной почты не корректный\n")
		return
	}
	if !letterValid {
		fmt.Fprintf(w, "Ошибка! Пароль должен содержать строчные буквы\n")
	}
	if !specValid {
		fmt.Fprintf(w, "Ошибка! Пароль должен содержать спец.символ\n")
	}
	if !lenValid {
		fmt.Fprintf(w, "Ошибка! Пароль должен содержать не менее 8 символов\n")
	}
	if !numbersValid {
		fmt.Fprintf(w, "Ошибка! Пароль должен содержать цифры\n")
	}
	if !containValid {
		fmt.Fprintf(w, "Ошибка! Пароль должен содержать прописные буквы\n")
	}
	if !weakValid {
		fmt.Fprintf(w, "Предупреждение! Очень слабый пароль, придумайте другой\n")
	}

	// Адрес электронной почты и пароль валидны
	if emailValid && letterValid && specValid && lenValid && numbersValid && containValid && weakValid {
		hash := check.HashPass(password)
		// Адрес электронной почты и пароль валидны
		c := storage.Account{
			Username: username,
			Password: hash,
		}

		// Проверяем есть ли такой пользователь в базе redis
		keys, err := api.db.KeysAccount(c)
		if err != nil {
			log.Println(err)
		}
		if keys == true {
			//  Если такой пользователь существует, отображаем сообщение об ошибке.
			fmt.Fprintf(w, "Такой пользователь уже существует")
			return
		}
		err = api.db.AddAccount(c)
		if err != nil {
			log.Println(err)
		}
		fmt.Fprintf(w, "Ваш аккаунт успешно создан.")
		http.Redirect(w, r, "/login", http.StatusFound)
	}
}

// Функция-обработчик для страницы с авторизацией
func (api *API) loginHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/login" {
		http.NotFound(w, r)
	}
	// Получаем данные из формы регистрации
	username := r.FormValue("username")
	password := r.FormValue("password")

	hash := check.HashPass(password)

	c := storage.Account{
		Username: username,
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
		log.Println(session.Name, session.Value, session.Secure)
		http.SetCookie(w, session)

		// Перенаправляем пользователя на защищенную страницу
		http.Redirect(w, r, "/dashboard", http.StatusFound)
	} else {
		// Если авторизация не удалась, отображаем сообщение об ошибке
		fmt.Fprintf(w, "Нет такой записи, проверте логин или пароль")
	}
}

// Функция-обработчик для защищенной страницы
func (api *API) dashboardHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/dashboard" {
		http.NotFound(w, r)
	}
	// Проверяем, авторизован ли пользователь
	session, err := r.Cookie("session")
	if err != nil || session.Value != "authenticated" {
		// Если пользователь не авторизован, перенаправляем на страницу входа
		http.Redirect(w, r, "/login", http.StatusFound)
		return
	}

	// Если пользователь авторизован, отображаем защищенную страницу
	fmt.Fprintf(w, "Добро пожаловать в панель управления!")
}

func (api *API) delAccountHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/delaccount" {
		http.NotFound(w, r)
	}
	// Получаем данные из формы регистрации
	username := r.FormValue("username")
	//password := r.FormValue("password")

	// Получаем логин из формы входа
	c := storage.Account{
		Username: username,
		//Password: password,
	}

	// Проверяем есть ли такой пользователь в базе redis
	keys, err := api.db.KeysAccount(c)
	if err != nil {
		log.Println(err)
	}
	if keys == true {
		a, err := api.db.DelAccount(c)
		if err != nil {
			log.Println(err)
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

			fmt.Fprintf(w, "Ваш аккаунт успешно удален.")
		}
	}
	if keys == false {
		//  Такой пользователь не существует, проверьте логин и пароль.
		fmt.Fprintf(w, "Такой пользователь не существует, проверьте логин.")
	}
}
