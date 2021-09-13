package comments

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/uid"
)

type PostgresComment struct {
	Uuid       uid.UID      `db:"uuid"`
	UserUuid   uid.UID      `db:"user_uuid"`
	SourceUuid uid.UID      `db:"source_uuid"`
	ParentUuid uid.UID      `db:"parent_uuid"`
	Body       string       `db:"body"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres *Postgres) Many(sourceId uid.UID) ([]Comment, error) {
	var comments []PostgresComment

	query := `SELECT * FROM comments WHERE source_uuid = $1 ORDER BY created_at`
	if err := postgres.Db.Select(&comments, query, sourceId); err != nil {
		return nil, err
	}

	return prepareMany(comments), nil
}

func (postgres *Postgres) Insert(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	var uuid uid.UID

	query := `
	INSERT INTO comments (user_uuid, source_uuid, parent_uuid, body) 
	VALUES ($1, $2, $3, $4) 
	RETURNING uuid`

	if err := postgres.Db.Get(&uuid, query, userId, sourceId, parentId, fields.Body); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres *Postgres) Update(commentId uid.UID, fields CommentFields) error {
	query := `UPDATE comments SET body = $2 WHERE uuid = $1`
	result, err := postgres.Db.Exec(query, commentId, fields.Body)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New(fmt.Sprintf("no updates were made to comment id [%s]", commentId))
	}

	return nil
}

func prepareOne(pc PostgresComment) Comment {
	return Comment{
		Id:        pc.Uuid,
		UserId:    pc.UserUuid,
		SourceId:  pc.SourceUuid,
		ParentId:  pc.ParentUuid,
		Body:      pc.Body,
		CreatedAt: pc.CreatedAt,
		UpdatedAt: pc.UpdatedAt.Time,
	}
}

func prepareMany(pc []PostgresComment) []Comment {
	var comments = make([]Comment, 0)

	for _, comment := range pc {
		comments = append(comments, prepareOne(comment))
	}

	return comments
}

func NewStorage(db *sqlx.DB) *Postgres {
	return &Postgres{db}
}
