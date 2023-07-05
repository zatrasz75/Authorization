package postgres

import (
	Interface "authorization/pkg/storage"
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"log"
	"time"
)

// Store Хранилище данных
type Store struct {
	db *pgxpool.Pool
}

// New Конструктор
func New(connstr string) (*Store, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	db, err := pgxpool.Connect(ctx, connstr)
	if err != nil {
		return nil, err
	}
	s := Store{
		db: db,
	}
	return &s, nil
}

// AddAccount Добавляет данные в базу Postgres
func (s *Store) AddAccount(c Interface.Account) error {
	_, err := s.db.Exec(context.Background(),
		"INSERT INTO accounts(username, password) VALUES ($1, $2);", c.Username, c.Password)
	if err != nil {
		log.Println(err)
		return err
	}

	return nil
}

// SearchAccount Находит пароль по ключу в базе Postgres
func (s *Store) SearchAccount(c Interface.Account) (string, error) {
	query := "SELECT password FROM accounts WHERE username = $1"

	var password string
	err := s.db.QueryRow(context.Background(), query, c.Username).Scan(&password)
	if err != nil {
		if err == pgx.ErrNoRows {
			return "", nil
		}
		return "", err
	}

	return password, nil
}

// KeysAccount Проверяет логин по ключу в базе Postgres
func (s *Store) KeysAccount(c Interface.Account) (bool, error) {
	query := "SELECT EXISTS(SELECT 1 FROM accounts WHERE username = $1)"

	var exists bool
	err := s.db.QueryRow(context.Background(), query, c.Username).Scan(&exists)
	if err != nil {
		return false, err
	}

	return exists, nil
}

// DelAccount Удаляет аккаунт в базе Postgres
func (s *Store) DelAccount(c Interface.Account) (bool, error) {
	delet := "DELETE FROM accounts WHERE username = $1"

	_, err := s.db.Exec(context.Background(), delet, c.Username)
	if err != nil {
		return false, err
	}

	return true, nil
}

// CreateAccountsTable Создает таблицу accounts
func (s *Store) CreateAccountsTable() error {
	qwery := `CREATE TABLE IF NOT EXISTS "accounts" (
    id SERIAL PRIMARY KEY,
    username TEXT NOT NULL,
    password TEXT NOT NULL
);`

	_, err := s.db.Exec(context.Background(), qwery)
	if err != nil {
		return err
	}

	return nil
}

// DropAccountsTable Удаляет таблицу accounts
func (s *Store) DropAccountsTable() error {
	drop := `DROP TABLE IF EXISTS "accounts";`

	_, err := s.db.Exec(context.Background(), drop)
	if err != nil {
		return err
	}

	return nil
}
