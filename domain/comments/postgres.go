package comments

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"quotes/pkg/uid"
)

type PostgresComment struct {
	Uuid       uid.UID      `db:"uuid"`
	UserUuid   uid.UID      `db:"user_uuid"`
	PostUuid   uid.UID      `db:"post_uuid"`
	ParentUuid uid.UID      `db:"parent_uuid"`
	Content    string       `db:"content"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres *Postgres) Many(postId uid.UID) ([]Comment, error) {
	var comments []PostgresComment

	query := "SELECT uuid, user_uuid, post_uuid, parent_uuid, content, created_at, updated_at, deleted_at FROM comments WHERE post_uuid = $1 ORDER BY created_at"
	if err := postgres.Db.Select(&comments, query, postId); err != nil {
		return nil, err
	}

	return prepareMany(comments), nil
}

func (postgres *Postgres) Insert(userId uid.UID, postId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	var uuid uid.UID

	query := "INSERT INTO comments (user_uuid, post_uuid, parent_uuid, content) VALUES ($1, $2, $3, $4) RETURNING uuid"
	if err := postgres.Db.Get(&uuid, query, userId, postId, parentId, fields.Content); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres *Postgres) Update(commentId uid.UID, fields CommentFields) (uid.UID, error) {
	var uuid uid.UID

	query := `UPDATE comments SET content = $2 WHERE uuid = $1 RETURNING uuid`
	if err := postgres.Db.Get(&uuid, query, commentId, fields.Content); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func prepareOne(pc PostgresComment) Comment {
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
