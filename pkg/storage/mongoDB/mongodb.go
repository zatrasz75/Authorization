package mongoDB

import (
	Interface "authorization/pkg/storage"
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type Storage struct {
	db *mongo.Client
}

const (
	databaseName   = "Account"  // имя БД
	collectionName = "accounts" // имя коллекции в БД
)

// New Конструктор, принимает строку подключения к БД.
func New(constr string) (*Storage, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	mongoOpts := options.Client().ApplyURI(constr).SetDirect(false)

	client, err := mongo.Connect(ctx, mongoOpts)
	if err != nil {
		return nil, err // возвращаем ошибку вместо использования log.Fatal
	}

	// Проверка подключения
	err = client.Ping(ctx, nil)
	if err != nil {
		return nil, err // возвращаем ошибку, если не удалось установить соединение
	}

	// не забываем закрывать ресурсы
	s := Storage{
		db: client,
	}
	return &s, nil
}

// AddAccount Добавляет данные в базу MongoDB
func (m *Storage) AddAccount(c Interface.Account) error {
	login := m.db.Database(databaseName).Collection(collectionName)
	_, err := login.InsertOne(context.Background(), c)
	if err != nil {
		return err
	}
	return nil
}

// SearchAccount Находит пароль по ключу в базе MongoDB
func (m *Storage) SearchAccount(c Interface.Account) (string, error) {
	// Получение коллекции accounts
	collection := m.db.Database(databaseName).Collection(collectionName)

	// Создание фильтра для поиска по ключу
	filter := bson.D{{"username", c.Username}}

	// Поиск документа по ключу
	var result Interface.Account
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Документ не найден
			return "", nil
		}
		// Возникла ошибка при выполнении запроса
		return "", err
	}
	// Возврат пароля из найденного документа
	return result.Password, nil
}

// KeysAccount Проверяет логин по ключу в базе MongoDB
func (m *Storage) KeysAccount(c Interface.Account) (bool, error) {
	// Получение коллекции accounts
	collection := m.db.Database(databaseName).Collection(collectionName)

	// Создание фильтра для поиска по ключу
	filter := bson.D{{"username", c.Username}}

	// Поиск документа по ключу
	var result Interface.Account
	err := collection.FindOne(context.Background(), filter).Decode(&result)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// Документ не найден
			return false, nil
		}
		// Возникла ошибка при выполнении запроса
		return false, err
	}

	// Документ найден
	return true, nil
}

// DelAccount Удаляет аккаунт в базе MongoDB
func (m *Storage) DelAccount(c Interface.Account) (bool, error) {
	// Получение коллекции accounts
	collection := m.db.Database(databaseName).Collection(collectionName)

	// Создание фильтра для поиска по ключу
	filter := bson.D{{"username", c.Username}}

	// Удаление документа по фильтру
	result, err := collection.DeleteOne(context.Background(), filter)
	if err != nil {
		// Возникла ошибка при выполнении запроса
		return false, err
	}

	// Проверка, был ли удален хотя бы один документ
	if result.DeletedCount > 0 {
		// Успешно удален хотя бы один документ
		return true, nil
	}

	// Если удаленных документов нет, возвращаем false
	return false, nil
}
