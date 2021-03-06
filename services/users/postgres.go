package users

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/uid"
)

type PostgresUser struct {
	Uuid           uid.UID        `db:"uuid"`
	AccountUuid    uid.UID        `db:"account_uuid"`
	Email          sql.NullString `db:"email"`
	Nickname       string         `db:"nickname"`
	ProfilePicture sql.NullString `db:"profile_picture"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}

type Postgres struct {
	db *sqlx.DB
}

func (p Postgres) Insert(accountId uid.UID, f *Fields) error {
	query := `INSERT INTO users (account_uuid, email, nickname, profile_picture) VALUES ($1, $2, $3, $4)`
	if _, err := p.db.Exec(query, accountId, f.Email, f.Nickname, f.ProfilePicture); err != nil {
		return err
	}

	return nil
}

func (p Postgres) ById(userId uid.UID) (User, error) {
	var user PostgresUser

	query := `SELECT * FROM users WHERE uuid = $1 LIMIT 1`
	if err := p.db.Get(&user, query, userId); err != nil {
		return User{}, err
	}

	return p.prepareOne(user), nil
}

func (p Postgres) ByIds(userIds []uid.UID) ([]User, error) {
	var users []PostgresUser

	query := `SELECT * FROM users WHERE uuid IN (?)`
	query, args, err := sqlx.In(query, userIds)
	if err != nil {
		return []User{}, err
	}
	query = p.db.Rebind(query)

	if err = p.db.Select(&users, query, args...); err != nil {
		return []User{}, err
	}

	return p.prepareMany(users), nil
}

func (p Postgres) ByAccountId(accountId uid.UID) (User, error) {
	var user PostgresUser

	query := `SELECT * FROM users WHERE account_uuid = $1 LIMIT 1`
	if err := p.db.Get(&user, query, accountId); err != nil {
		return User{}, err
	}

	return p.prepareOne(user), nil
}

func (p Postgres) prepareMany(pu []PostgresUser) []User {
	var users = make([]User, 0)

	for _, post := range pu {
		users = append(users, p.prepareOne(post))
	}

	return users
}

func (Postgres) prepareOne(pu PostgresUser) User {
	return User{
		Id:             pu.Uuid,
		Email:          pu.Email.String,
		Nickname:       pu.Nickname,
		ProfilePicture: pu.ProfilePicture.String,
		CreatedAt:      pu.CreatedAt,
		UpdatedAt:      pu.UpdatedAt.Time,
	}
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
