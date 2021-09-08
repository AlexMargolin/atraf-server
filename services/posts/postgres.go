package posts

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/uid"
)

type PostgresPost struct {
	Uuid      uid.UID      `db:"uuid"`
	UserUuid  uid.UID      `db:"user_uuid"`
	Body      string       `db:"body"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres Postgres) One(postId uid.UID) (Post, error) {
	var post PostgresPost

	query := "SELECT uuid, user_uuid, body, created_at, updated_at, deleted_at FROM posts WHERE uuid = $1 LIMIT 1"
	if err := postgres.Db.Get(&post, query, postId); err != nil {
		return Post{}, err
	}

	return prepareOne(post), nil
}

func (postgres Postgres) Many(p *middleware.PaginationContext) ([]Post, error) {
	var posts []PostgresPost

	if p.Cursor.Key != uid.Nil {
		query := "SELECT uuid, user_uuid, body, created_at, updated_at, deleted_at FROM posts WHERE uuid < $1 AND created_at < $2 ORDER BY created_at DESC LIMIT $3"
		if err := postgres.Db.Select(&posts, query, p.Cursor.Key, p.Cursor.Value, p.Limit); err != nil {
			return nil, err
		}
	} else {
		query := "SELECT uuid, user_uuid, body, created_at, updated_at, deleted_at FROM posts ORDER BY created_at DESC LIMIT $1"
		if err := postgres.Db.Select(&posts, query, p.Limit); err != nil {
			return nil, err
		}
	}

	return prepareMany(posts), nil
}

func (postgres Postgres) Insert(userId uid.UID, fields PostFields) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO posts (user_uuid, body) VALUES ($1, $2) RETURNING uuid"
	if err := postgres.Db.Get(&uuid, query, userId, fields.Body); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres Postgres) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	var uuid uid.UID

	query := "UPDATE posts SET body = $2 WHERE uuid = $1 RETURNING uuid"
	if err := postgres.Db.Get(&uuid, query, postId, fields.Body); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func prepareOne(pp PostgresPost) Post {
	return Post{
		Id:        pp.Uuid,
		UserId:    pp.UserUuid,
		Body:      pp.Body,
		CreatedAt: pp.CreatedAt,
		UpdatedAt: pp.UpdatedAt.Time,
	}
}

func prepareMany(pp []PostgresPost) []Post {
	var posts = make([]Post, 0)

	for _, post := range pp {
		posts = append(posts, prepareOne(post))
	}

	return posts
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
