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
	FirstName      sql.NullString `db:"first_name"`
	LastName       sql.NullString `db:"last_name"`
	ProfilePicture sql.NullString `db:"profile_picture"`
	CreatedAt      time.Time      `db:"created_at"`
	UpdatedAt      sql.NullTime   `db:"updated_at"`
	DeletedAt      sql.NullTime   `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres Postgres) One(userId uid.UID) (User, error) {
	var user PostgresUser

	query := "SELECT uuid, account_uuid, email, first_name, last_name, profile_picture, created_at, updated_at, deleted_at FROM users WHERE uuid = $1 LIMIT 1"
	if err := postgres.Db.Get(&user, query, userId); err != nil {
		return User{}, err
	}

	return prepareOne(user), nil
}

func (postgres Postgres) Insert(accountId uid.UID, fields UserFields) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO users (account_uuid, email, first_name, last_name, profile_picture) VALUES ($1, $2, $3, $4, $5) RETURNING uuid"
	if err := postgres.Db.Get(&uuid, query, accountId, fields.Email, fields.FirstName, fields.LastName, fields.ProfilePicture); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func prepareOne(pu PostgresUser) User {
	return User{
		Id:             pu.Uuid,
		Email:          pu.Email.String,
		FirstName:      pu.FirstName.String,
		LastName:       pu.LastName.String,
		ProfilePicture: pu.ProfilePicture.String,
	}
}

func prepareMany(pp []PostgresUser) []User {
	var users = make([]User, 0)

	for _, post := range pp {
		users = append(users, prepareOne(post))
	}

	return users
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
