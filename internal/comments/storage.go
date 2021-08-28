package comments

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"quotes/pkg/uid"
)

type SqlComment struct {
	Id        uid.UID
	UserId    uid.UID
	PostId    uid.UID
	ParentId  uid.UID
	Content   string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
	DeletedAt sql.NullTime
}

type SqlStorage struct {
	Db *sql.DB
}

func (storage *SqlStorage) Many(postId uid.UID) ([]Comment, error) {
	comments := make([]Comment, 0)

	query := "SELECT * FROM comments WHERE post_id = ? ORDER BY created_at"
	rows, err := storage.Db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan Rows
	for rows.Next() {
		var s SqlComment

		err = rows.Scan(&s.Id, &s.UserId, &s.PostId, &s.ParentId, &s.Content, &s.CreatedAt, &s.UpdatedAt, &s.DeletedAt)
		if err != nil {
			return nil, err
		}

		comments = append(comments, toComment(s))
	}

	return comments, nil
}

func (storage *SqlStorage) Insert(userId uid.UID, postId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	commentId := uid.New()

	query := "INSERT INTO comments (id, user_id, post_id, parent_id, content) VALUES (?, ?, ?, ?, ?)"
	if _, err := storage.Db.Exec(query, commentId, userId, postId, parentId, fields.Content); err != nil {
		return uid.Nil, err
	}

	return commentId, nil
}

func (storage *SqlStorage) Update(commentId uid.UID, fields CommentFields) (uid.UID, error) {
	query := "UPDATE comments SET content = ? WHERE id = ? LIMIT 1"

	result, err := storage.Db.Exec(query, fields.Content, commentId)
	if err != nil {
		return uid.Nil, err
	}

	ra, err := result.RowsAffected()
	if err != nil {
		return uid.Nil, err
	}

	if ra == 0 {
		return uid.Nil, errors.New(fmt.Sprintf("unable to update comment with the id [%s]", commentId))
	}

	return commentId, nil
}

func toComment(sqlComment SqlComment) Comment {
	return Comment{
		Id:        sqlComment.Id,
		UserId:    sqlComment.UserId,
		PostId:    sqlComment.PostId,
		ParentId:  sqlComment.ParentId,
		Content:   sqlComment.Content,
		CreatedAt: sqlComment.CreatedAt,
		UpdatedAt: sqlComment.UpdatedAt.Time,
	}
}

func NewStorage(db *sql.DB) *SqlStorage {
	return &SqlStorage{db}
}
