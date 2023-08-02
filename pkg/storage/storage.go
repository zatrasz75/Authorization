package storage

type Account struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type FormAccount struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type Response struct {
	Success       bool     `json:"success"`
	Message       string   `json:"message"`
	ErrorMessages []string `json:"errorMessages"`
}

type Interface interface {
	AddAccount(c Account) error
	SearchAccount(c Account) (string, error)
	KeysAccount(c Account) (bool, error)
	DelAccount(c Account) (bool, error)
}
