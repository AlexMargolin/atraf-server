package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/uid"
)

type PostgresAccount struct {
	Uuid         uid.UID      `db:"uuid"`
	Email        string       `db:"email"`
	PasswordHash []byte       `db:"password_hash"`
	Active       bool         `db:"active"`
	CreatedAt    time.Time    `db:"created_at"`
	UpdatedAt    sql.NullTime `db:"updated_at"`
	DeletedAt    sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres *Postgres) NewAccount(email string, passwordHash []byte) (Account, error) {
	var account PostgresAccount

	query := "INSERT INTO accounts (email, password_hash) VALUES ($1, $2) RETURNING *"
	if err := postgres.Db.Get(&account, query, email, passwordHash); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (postgres *Postgres) ByAccountId(accountId uid.UID) (Account, error) {
	var account PostgresAccount

	query := "SELECT * FROM accounts WHERE uuid = $1 LIMIT 1"
	if err := postgres.Db.Get(&account, query, accountId); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (postgres *Postgres) ByEmail(email string) (Account, error) {
	var account PostgresAccount

	query := "SELECT * FROM accounts WHERE email = $1 LIMIT 1"
	if err := postgres.Db.Get(&account, query, email); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (postgres *Postgres) UpdatePassword(accountId uid.UID, passwordHash []byte) error {
	query := "UPDATE accounts SET password_hash = $2 WHERE uuid = $1"

	result, err := postgres.Db.Exec(query, accountId, passwordHash)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("password couldn't be updated")
	}

	return nil
}

func (postgres *Postgres) UpdateStatus(accountId uid.UID, active bool) error {
	query := "UPDATE accounts SET active = $2 WHERE uuid = $1"

	result, err := postgres.Db.Exec(query, accountId, active)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("status couldn't be updated")
	}

	return nil
}

func prepareOne(pa PostgresAccount) Account {
	return Account{
		Id:           pa.Uuid,
		Email:        pa.Email,
		PasswordHash: pa.PasswordHash,
		Active:       pa.Active,
		CreatedAt:    pa.CreatedAt,
		UpdatedAt:    pa.UpdatedAt.Time,
	}
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
