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
	db *sqlx.DB
}

func (p Postgres) Insert(userId uid.UID, sourceId uid.UID, parentId uid.UID, f *Fields) (Comment, error) {
	var c PostgresComment

	query := `
	INSERT INTO comments (user_uuid, source_uuid, parent_uuid, body) 
	VALUES ($1, $2, $3, $4) 
	RETURNING *`

	if err := p.db.Get(&c, query, userId, sourceId, parentId, f.Body); err != nil {
		return Comment{}, err
	}

	return prepareOne(c), nil
}

func (p Postgres) Update(commentId uid.UID, f *Fields) error {
	query := `UPDATE comments SET body = $2 WHERE uuid = $1`
	result, err := p.db.Exec(query, commentId, f.Body)
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

func (p Postgres) Many(sourceId uid.UID) ([]Comment, error) {
	var c []PostgresComment

	query := `SELECT * FROM comments WHERE source_uuid = $1 ORDER BY created_at DESC`
	if err := p.db.Select(&c, query, sourceId); err != nil {
		return nil, err
	}

	return prepareMany(c), nil
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
