package account

import (
	"database/sql"

	"quotes/pkg/uid"
)

type SqlAccount struct {
	Id           uid.UID
	Email        string
	PasswordHash []byte
	CreatedAt    sql.NullTime
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
}

type SqlStorage struct {
	Db *sql.DB
}

// Insert create a new user in the database
func (storage *SqlStorage) Insert(email string, passwordHash []byte) (uid.UID, error) {
	accountId := uid.New()

	query := "INSERT INTO accounts (id, email, password_hash) VALUES (?, ?, ?)"
	if _, err := storage.Db.Exec(query, accountId, email, passwordHash); err != nil {
		return uid.Nil, err
	}

	return accountId, nil
}

// ByEmail retrieves an account by its email address.
// Returns an error when no account is found
func (storage *SqlStorage) ByEmail(email string) (Account, error) {
	var s SqlAccount

	query := "SELECT * FROM accounts WHERE email = ? LIMIT 1"
	row := storage.Db.QueryRow(query, email)

	err := row.Scan(&s.Id, &s.Email, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		return Account{}, err
	}

	return storage.toAccount(s), nil
}

// toAccount converts an SqlAccount into Account
func (SqlStorage) toAccount(s SqlAccount) Account {
	return Account{
		Id:           s.Id,
		Email:        s.Email,
		PasswordHash: s.PasswordHash,
		CreatedAt:    s.CreatedAt.Time,
		UpdatedAt:    s.UpdatedAt.Time,
		DeletedAt:    s.DeletedAt.Time,
	}
}

// NewStorage returns a new MySQL Storage instance
func NewStorage(db *sql.DB) *SqlStorage {
	return &SqlStorage{db}
}
