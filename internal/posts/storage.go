package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"quotes/pkg/uid"
)

type SqlPost struct {
	Id        uid.UID
	UserId    uid.UID
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type SqlStorage struct {
	Db *sql.DB
}

func (storage *SqlStorage) Count() (int, error) {
	var count int

	query := "SELECT COUNT(*) FROM posts"
	row := storage.Db.QueryRow(query)

	err := row.Scan(&count)
	if err != nil {
		return count, err
	}

	return count, nil
}

func (storage *SqlStorage) One(postId uid.UID) (Post, error) {
	var s SqlPost

	query := "SELECT * FROM posts WHERE id = ? LIMIT 1"
	row := storage.Db.QueryRow(query, postId)

	err := row.Scan(&s.Id, &s.UserId, &s.Content, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
	if err != nil {
		return Post{}, err
	}

	return storage.toPost(s), nil
}

func (storage *SqlStorage) Many(offset int, limit int) ([]Post, error) {
	posts := make([]Post, 0)

	query := "SELECT * FROM posts ORDER BY created_at DESC LIMIT ? OFFSET ?"
	rows, err := storage.Db.Query(query, limit, offset)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	// Scan Rows
	for rows.Next() {
		var s SqlPost

		err = rows.Scan(&s.Id, &s.UserId, &s.Content, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
		if err != nil {
			return nil, err
		}

		posts = append(posts, storage.toPost(s))
	}

	return posts, nil
}

func (storage *SqlStorage) Insert(userId uid.UID, fields PostFields) (uid.UID, error) {
	postId := uid.New()

	query := "INSERT INTO posts (id, user_id, content) VALUES (?, ?, ?)"
	if _, err := storage.Db.Exec(query, postId, userId, fields.Content); err != nil {
		return uid.Nil, err
	}

	return postId, nil
}

func (storage *SqlStorage) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	query := "UPDATE posts SET content = ? WHERE id = ? LIMIT 1"

	result, err := storage.Db.Exec(query, fields.Content, postId)
	if err != nil {
		return uid.Nil, err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return uid.Nil, err
	}

	if ra == 0 {
		return uid.Nil, errors.New(fmt.Sprintf("unable to update post with the id [%s]", postId))
	}

	return postId, nil
}

// Receives a SqlPost and returns a Post
func (SqlStorage) toPost(sqlPost SqlPost) Post {
	return Post{
		Id:        sqlPost.Id,
		UserId:    sqlPost.UserId,
		Content:   sqlPost.Content,
		CreatedAt: sqlPost.CreatedAt,
		UpdatedAt: sqlPost.UpdatedAt.Time,
	}
}

// NewStorage returns a new MySQL Storage instance
func NewStorage(db *sql.DB) *SqlStorage {
	return &SqlStorage{db}
}
