package comments

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"quotes/pkg/uid"
)

type PostgresComment struct {
	Uuid       uid.UID
	UserUuid   uid.UID
	PostUuid   uid.UID
	ParentUuid uid.UID
	Content    string
	CreatedAt  time.Time
	UpdatedAt  sql.NullTime
	DeletedAt  sql.NullTime
}

type Postgres struct {
	Db *sql.DB
}

func (postgres *Postgres) Many(postId uid.UID) ([]Comment, error) {
	comments := make([]Comment, 0)

	query := "SELECT uuid, user_uuid, post_uuid, parent_uuid, content, created_at, updated_at, deleted_at FROM comments WHERE post_uuid = $1 ORDER BY created_at"

	rows, err := postgres.Db.Query(query, postId)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	// Scan Rows
	for rows.Next() {
		var pc PostgresComment

		err = rows.Scan(
			&pc.Uuid,
			&pc.UserUuid,
			&pc.PostUuid,
			&pc.ParentUuid,
			&pc.Content,
			&pc.CreatedAt,
			&pc.UpdatedAt,
			&pc.DeletedAt,
		)
		if err != nil {
			return nil, err
		}

		comments = append(comments, toComment(pc))
	}

	return comments, nil
}

func (postgres *Postgres) Insert(userId uid.UID, postId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO comments (user_uuid, post_uuid, parent_uuid, content) VALUES ($1, $2, $3, $4) RETURNING uuid"

	err := postgres.Db.QueryRow(query, userId, postId, parentId, fields.Content).Scan(
		&uuid,
	)
	if err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres *Postgres) Update(commentId uid.UID, fields CommentFields) (uid.UID, error) {
	query := `UPDATE comments SET content = $2 WHERE uuid = $1`

	result, err := postgres.Db.Exec(query, commentId, fields.Content)
	if err != nil {
		return uid.Nil, err
	}

	if rows, err := result.RowsAffected(); err != nil || rows == 0 {
		return uid.Nil, errors.New(fmt.Sprintf("0 rows affected when updating Post[%s]", commentId))
	}

	return commentId, nil
}

func toComment(pc PostgresComment) Comment {
	return Comment{
		Id:        pc.Uuid,
		UserId:    pc.UserUuid,
		PostId:    pc.PostUuid,
		ParentId:  pc.ParentUuid,
		Content:   pc.Content,
		CreatedAt: pc.CreatedAt,
		UpdatedAt: pc.UpdatedAt.Time,
	}
}

func NewStorage(db *sql.DB) *Postgres {
	return &Postgres{db}
}
