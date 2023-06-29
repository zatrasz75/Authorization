package check

import (
	"fmt"
	"log"
	"regexp"
	"strings"
)

const (
	minLength = 8
)

// CheckPassword Функция для проверки валидности пароля
func CheckPassword(password string) bool {

	// Проверяем, что пароль имеет длину не менее 8 символов
	lenRegex := regexp.MustCompile(fmt.Sprintf(`^.{%d,}$`, minLength))
	if !lenRegex.MatchString(password) {
		log.Printf("Ошибка! Длина пароля менее %d\n", minLength)
		return false
	}

	// Проверяем, что пароль содержит хотя бы одну цифру
	if !regexp.MustCompile(`[0-9]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать цифры")
		return false
	}

	// Проверяем, что пароль содержит хотя бы одну заглавную букву
	if !regexp.MustCompile(`[A-Z]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать прописные буквы")
		return false
	}

	// Проверяем, что пароль содержит хотя бы одну строчную букву
	if !regexp.MustCompile(`[a-z]+`).MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать строчные буквы")
		return false
	}

	// Проверяем что пароль должен содержать спец.символ
	specCharRegex := regexp.MustCompile(`[!@#$%^&*()\-=+,./\\_]+`)
	if !specCharRegex.MatchString(password) {
		log.Println("Ошибка! Пароль должен содержать спец.символ")
		return false
	}

	// Проверяем что пароль не слабый
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
