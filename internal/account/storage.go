package account

import (
	"database/sql"
	"time"

	"quotes/pkg/uid"
)

type SqlAccount struct {
	Uuid         uid.UID
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
}

type SqlStorage struct {
	Db *sql.DB
}

func (storage *SqlStorage) Insert(email string, passwordHash []byte) (uid.UID, error) {
	accountId := uid.New()

	query := "INSERT INTO accounts (uuid, email, password_hash) VALUES (?, ?, ?)"
	if _, err := storage.Db.Exec(query, accountId, email, passwordHash); err != nil {
		return uid.Nil, err
	}

	return accountId, nil
}

func (storage *SqlStorage) ByEmail(email string) (Account, error) {
	var s SqlAccount

	query := "SELECT * FROM accounts WHERE email = ? LIMIT 1"
	row := storage.Db.QueryRow(query, email)

	err := row.Scan(&s.Uuid, &s.Email, &s.PasswordHash, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		return Account{}, err
	}

	return storage.toAccount(s), nil
}

func (SqlStorage) toAccount(s SqlAccount) Account {
	return Account{
		Id:           s.Uuid,
		Email:        s.Email,
		PasswordHash: s.PasswordHash,
		CreatedAt:    s.CreatedAt,
		UpdatedAt:    s.UpdatedAt.Time,
	}
}

func NewStorage(db *sql.DB) *SqlStorage {
	return &SqlStorage{db}
}
