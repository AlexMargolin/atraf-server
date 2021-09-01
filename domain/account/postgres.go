package account

import (
	"database/sql"
	"time"

	"quotes/pkg/uid"
)

type PostgresAccount struct {
	Uuid         uid.UID
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    sql.NullTime
	DeletedAt    sql.NullTime
}

type Postgres struct {
	Db *sql.DB
}

func (postgres *Postgres) Insert(email string, passwordHash []byte) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO accounts (email, password_hash) VALUES ($1, $2) RETURNING uuid"

	err := postgres.Db.QueryRow(query, email, passwordHash).Scan(
		&uuid,
	)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres *Postgres) ByEmail(email string) (Account, error) {
	var pa PostgresAccount

	query := "SELECT uuid, email, password_hash, created_at, updated_at, deleted_at FROM accounts WHERE email = $1 LIMIT 1"

	err := postgres.Db.QueryRow(query, email).Scan(
		&pa.Uuid,
		&pa.Email,
		&pa.PasswordHash,
		&pa.CreatedAt,
		&pa.UpdatedAt,
		&pa.DeletedAt,
	)
	if err != nil {
		return Account{}, err
	}

	return prepare(pa), nil
}

func prepare(pa PostgresAccount) Account {
	return Account{
		Id:           pa.Uuid,
		Email:        pa.Email,
		PasswordHash: pa.PasswordHash,
		CreatedAt:    pa.CreatedAt,
		UpdatedAt:    pa.UpdatedAt.Time,
	}
}

func NewStorage(db *sql.DB) *Postgres {
	return &Postgres{db}
}
