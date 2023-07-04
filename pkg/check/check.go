package check

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	minLength = 8
)

// CheckEmail Функция для проверки валидности адреса электронной почты
func CheckEmail(email string) bool {
	// Регулярное выражение для проверки адреса электронной почты
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	if !regex.MatchString(email) {
		log.Println("Не корректный адрес")
		return false
	}
	return true
}

// LowercaseLetter Проверяем, что пароль содержит хотя бы одну строчную букву
func LowercaseLetter(password string) bool {
	if !regexp.MustCompile(`[a-z]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать строчные буквы")
		return false
	}
	return true
}

// SpecCharRegex Проверяем что пароль должен содержать спец. Символ
func SpecCharRegex(password string) bool {
	specCharRegex := regexp.MustCompile(`[!@#$%^&*()\-=+,./\\_]+`)
	if !specCharRegex.MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать спец.символ")
		return false
	}
	return true
}

// LenPass Проверяем, что пароль имеет длину не менее 8 символов
func LenPass(password string) bool {
	lenRegex := regexp.MustCompile(fmt.Sprintf(`^.{%d,}$`, minLength))
	if !lenRegex.MatchString(password) {
		log.Printf("Ошибка! Длина пароля менее %d\n", minLength)
		return false
	}
	return true
}

// NumbersPass Проверяем, что пароль содержит хотя бы одну цифру
func NumbersPass(password string) bool {
	if !regexp.MustCompile(`[0-9]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать цифры")
		return false
	}
	return true
}

// ContainPass Проверяем, что пароль содержит хотя бы одну заглавную букву
func ContainPass(password string) bool {
	if !regexp.MustCompile(`[A-Z]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать прописные буквы")
		return false
	}
	return true
}

// WeakPass Проверяем что пароль не слабый
func WeakPass(password string) bool {

	mostPopularPassword := []string{
		"Qq123456",
		"Qwerty123",
	}
	join := strings.Join(mostPopularPassword, "|")
	weakPassRegex := regexp.MustCompile(fmt.Sprintf("^(%s)$", join))
	if weakPassRegex.MatchString(password) {
		log.Println("Предупреждение! Очень слабый пароль, придумайте другой")
		return false
	}

	return true
}

// HashPass Хеширование пароля.
func HashPass(password string) string {
	// Создание нового объекта хэша SHA-256
	hash := sha256.New()

	// Запись данных в хэш-функцию
	hash.Write([]byte(password))

	// Получение окончательного хэш-значения в виде среза байт
	hashBytes := hash.Sum(nil)

	// Преобразование хэш-значения в строку в шестнадцатеричном формате
	hashString := hex.EncodeToString(hashBytes)

	return hashString
}
