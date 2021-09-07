package account

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/uid"
)

type PostgresAccount struct {
	Uuid         uid.UID      `db:"uuid"`
	Email        string       `db:"email"`
	PasswordHash []byte       `db:"password_hash"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres *Postgres) Insert(email string, passwordHash []byte) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO accounts (email, password_hash) VALUES ($1, $2) RETURNING uuid"
	if err := postgres.Db.Get(&uuid, query, email, passwordHash); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres *Postgres) ByEmail(email string) (Account, error) {
	var account PostgresAccount

	query := "SELECT uuid, email, password_hash, created_at, updated_at, deleted_at FROM accounts WHERE email = $1 LIMIT 1"
	if err := postgres.Db.Get(&account, query, email); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func prepareOne(pa PostgresAccount) Account {
	return Account{
		Id:           pa.Uuid,
		Email:        pa.Email,
		PasswordHash: pa.PasswordHash,
		CreatedAt:    pa.CreatedAt,
		UpdatedAt:    pa.UpdatedAt.Time,
	}
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
