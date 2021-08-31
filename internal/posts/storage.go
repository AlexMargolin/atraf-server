package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"quotes/pkg/uid"
)

type PostgresPost struct {
	Uuid      uid.UID
	UserUuid  uid.UID
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type Postgres struct {
	Db *sql.DB
}

func (postgres Postgres) Count() (int, error) {
	var count int

	query := `SELECT COUNT(*) FROM posts`

	row := postgres.Db.QueryRow(query)
	if err := row.Scan(&count); err != nil {
		return count, err
	}

	return count, nil
}

func (postgres Postgres) One(postId uid.UID) (Post, error) {
	var pp PostgresPost

	query := `SELECT uuid, user_uuid, content, created_at, updated_at, deleted_at 
			  FROM posts 
			  WHERE uuid = $1 
			  LIMIT 1`

	row := postgres.Db.QueryRow(query, postId)
	err := row.Scan(
		&pp.Uuid,
		&pp.UserUuid,
		&pp.Content,
		&pp.CreatedAt,
		&pp.UpdatedAt,
		&pp.DeletedAt,
	)
	if err != nil {
		return Post{}, err
	}

	return postgres.toPost(pp), nil
}

func (postgres Postgres) Many(offset int, limit int) ([]Post, error) {
	posts := make([]Post, 0)

	query := `SELECT uuid, user_uuid, content, created_at, updated_at, deleted_at 
			  FROM posts	
			  ORDER BY created_at DESC 
			  OFFSET $1 LIMIT $2`

	rows, err := postgres.Db.Query(query, offset, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var s PostgresPost

		err = rows.Scan(
			&s.Uuid,
			&s.UserUuid,
			&s.Content,
			&s.CreatedAt,
			&s.UpdatedAt,
			&s.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		posts = append(posts, postgres.toPost(s))
	}

	return posts, nil
}

func (postgres Postgres) Insert(userId uid.UID, fields PostFields) (uid.UID, error) {
	var uuid uid.UID

	query := `INSERT INTO posts (user_uuid, content) 
			  VALUES ($1, $2) 
			  RETURNING uuid`

	row := postgres.Db.QueryRow(query, userId, fields.Content)

	err := row.Scan(&uuid)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres Postgres) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	query := `UPDATE posts 
			  SET content = $1 
			  WHERE uuid = $2`

	result, err := postgres.Db.Exec(query, fields.Content, postId)
	if err != nil {
		return uid.Nil, err
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return uid.Nil, errors.New(fmt.Sprintf("0 rows affected when updating Post[%s]", postId))
	}

	return postId, nil
}

// Receives a PostgresPost and returns a Post
func (Postgres) toPost(pp PostgresPost) Post {
	return Post{
		Id:        pp.Uuid,
		UserId:    pp.UserUuid,
		Content:   pp.Content,
		CreatedAt: pp.CreatedAt,
		UpdatedAt: pp.UpdatedAt.Time,
	}
}

// NewStorage returns a new MySQL Storage instance
func NewStorage(db *sql.DB) *Postgres {
	return &Postgres{db}
}
