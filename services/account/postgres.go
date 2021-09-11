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

func (postgres *Postgres) Insert(email string, passwordHash []byte) (uid.UID, int, error) {
	var code int
	var uuid uid.UID

	tx, err := postgres.Db.Beginx()
	if err != nil {
		return uuid, code, err
	}

	defer func() {
		if err != nil {
			_ = tx.Rollback()
		} else {
			_ = tx.Commit()
		}
	}()

	accountQuery := `
	INSERT INTO accounts (email, password_hash) 
		VALUES ($1, $2) 
	RETURNING uuid`

	activationsQuery := `
	INSERT INTO accounts_activations 
		VALUES ($1, floor(random() * (999999 - 100000 + 1) + 100000)) 
	RETURNING activation_code`

	if err = tx.Get(&uuid, accountQuery, email, passwordHash); err != nil {
		return uuid, code, err
	}

	if err = tx.Get(&code, activationsQuery, uuid); err != nil {
		return uuid, code, err
	}

	return uuid, code, nil
}

func (postgres *Postgres) ById(accountId uid.UID) (Account, error) {
	var account PostgresAccount

	query := `
	SELECT uuid, email, password_hash, active, created_at, updated_at, deleted_at 
	FROM accounts 
	WHERE uuid = $1
	LIMIT 1`

	if err := postgres.Db.Get(&account, query, accountId); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (postgres *Postgres) ByEmail(email string) (Account, error) {
	var account PostgresAccount

	query := `
	SELECT uuid, email, password_hash, active, created_at, updated_at, deleted_at 
	FROM accounts 
	WHERE email = $1
	LIMIT 1`

	if err := postgres.Db.Get(&account, query, email); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (postgres *Postgres) UpdatePassword(accountId uid.UID, passwordHash []byte) error {
	query := "UPDATE accounts SET password_hash = $1 WHERE uuid = $2"

	result, err := postgres.Db.Exec(query, passwordHash, accountId)
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
