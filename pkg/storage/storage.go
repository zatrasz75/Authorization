package storage

type Account struct {
	Username string
	Password string
}

type Interface interface {
	AddAccount(c Account) error
	SearchAccount(c Account) (string, error)
	KeysAccount(c Account) (bool, error)
	DelAccount(c Account) (bool, error)
	//DropAccountsTable(c Account) error
}
