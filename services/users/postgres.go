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
	db *sqlx.DB
}

func (p Postgres) ById(userId uid.UID) (User, error) {
	var user PostgresUser

	query := `
	SELECT uuid, 
	       account_uuid,
	       email,
	       first_name,
	       last_name,
	       profile_picture,
	       created_at,
	       updated_at,
	       deleted_at 
	FROM users 
	WHERE uuid = $1 
	LIMIT 1`

	if err := p.db.Get(&user, query, userId); err != nil {
		return User{}, err
	}

	return prepareOne(user), nil
}

func (p Postgres) ByIds(userIds []uid.UID) ([]User, error) {
	var users []PostgresUser

	query := `
	SELECT uuid,
	       account_uuid,
	       email,
	       first_name,
	       last_name,
	       profile_picture,
	       created_at,
	       updated_at,
	       deleted_at 
	FROM users 
	WHERE uuid IN (?)`

	query, args, err := sqlx.In(query, userIds)
	if err != nil {
		return []User{}, err
	}
	query = p.db.Rebind(query)

	if err := p.db.Select(&users, query, args...); err != nil {
		return []User{}, err
	}

	return prepareMany(users), nil
}

func (p Postgres) ByAccountId(accountId uid.UID) (User, error) {
	var user PostgresUser

	query := `
	SELECT uuid,
	       account_uuid,
	       email,
	       first_name,
	       last_name,
	       profile_picture,
	       created_at,
	       updated_at,
	       deleted_at 
	FROM users 
	WHERE account_uuid = $1 
	LIMIT 1`

	if err := p.db.Get(&user, query, accountId); err != nil {
		return User{}, err
	}

	return prepareOne(user), nil
}

func (p Postgres) Insert(accountId uid.UID, fields UserFields) error {
	query := `
	INSERT INTO users (account_uuid, email, first_name, last_name, profile_picture) 
	VALUES ($1, $2, $3, $4, $5)`

	if _, err := p.db.Exec(query, accountId, fields.Email, fields.FirstName, fields.LastName, fields.ProfilePicture); err != nil {
		return err
	}

	return nil
}

func prepareOne(pu PostgresUser) User {
	return User{
		Id:             pu.Uuid,
		Email:          pu.Email.String,
		FirstName:      pu.FirstName.String,
		LastName:       pu.LastName.String,
		ProfilePicture: pu.ProfilePicture.String,
		CreatedAt:      pu.CreatedAt,
	}
}

func prepareMany(pu []PostgresUser) []User {
	var users = make([]User, 0)

	for _, post := range pu {
		users = append(users, prepareOne(post))
	}

	return users
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
