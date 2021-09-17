package account

import (
	"database/sql"
	"errors"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/uid"
)

type PostgresAccount struct {
	Uuid           uid.UID        `db:"uuid"`
	Email          string         `db:"email"`
	PasswordHash   []byte         `db:"password_hash"`
	Active         bool           `db:"active"`
	ActivationCode sql.NullString `db:"activation_code"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}

type Postgres struct {
	db *sqlx.DB
}

func (p Postgres) NewAccount(email string, passwordHash []byte) (Account, error) {
	var account PostgresAccount

	query := `INSERT INTO accounts (email, password_hash) VALUES ($1, $2) RETURNING *`
	if err := p.db.Get(&account, query, email, passwordHash); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (p Postgres) SetPending(accountId uid.UID) (string, error) {
	var code string

	query := `
	UPDATE accounts 
	SET active = false, 
	    activation_code = DEFAULT
	WHERE uuid = $1
	RETURNING activation_code`

	if err := p.db.Get(&code, query, accountId); err != nil {
		return code, err
	}

	return code, nil
}

func (p Postgres) SetActive(accountId uid.UID, activationCode string) error {
	query := `
	UPDATE accounts 
	SET active = true, 
	    activation_code = NULL
	WHERE uuid = $1 
	  AND activation_code = $2`

	result, err := p.db.Exec(query, accountId, activationCode)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New("account couldn't be activated")
	}

	return nil
}

func (p Postgres) ByAccountId(accountId uid.UID) (Account, error) {
	var account PostgresAccount

	query := `SELECT * FROM accounts WHERE uuid = $1 LIMIT 1`
	if err := p.db.Get(&account, query, accountId); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (p Postgres) ByEmail(email string) (Account, error) {
	var account PostgresAccount

	query := `SELECT * FROM accounts WHERE email = $1 LIMIT 1`
	if err := p.db.Get(&account, query, email); err != nil {
		return Account{}, err
	}

	return prepareOne(account), nil
}

func (p Postgres) UpdatePassword(accountId uid.UID, passwordHash []byte) error {
	query := `UPDATE accounts SET password_hash = $2 WHERE uuid = $1`

	result, err := p.db.Exec(query, accountId, passwordHash)
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
		Id:             pa.Uuid,
		Email:          pa.Email,
		PasswordHash:   pa.PasswordHash,
		Active:         pa.Active,
		ActivationCode: pa.ActivationCode.String,
		CreatedAt:      pa.CreatedAt,
		UpdatedAt:      pa.UpdatedAt.Time,
	}
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
