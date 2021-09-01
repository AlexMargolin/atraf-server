package posts

import (
	"database/sql"
	"time"

	"github.com/jmoiron/sqlx"

	"quotes/pkg/uid"
)

type PostgresPost struct {
	Uuid      uid.UID      `db:"uuid"`
	UserUuid  uid.UID      `db:"user_uuid"`
	Content   string       `db:"content"`
	CreatedAt time.Time    `db:"created_at"`
	UpdatedAt sql.NullTime `db:"updated_at"`
	DeletedAt sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	Db *sqlx.DB
}

func (postgres Postgres) One(postId uid.UID) (Post, error) {
	var post PostgresPost

	query := `
	SELECT uuid,
	       user_uuid,
	       content,
	       created_at,
	       updated_at,
	       deleted_at 
	FROM posts 
	WHERE uuid = $1 
	LIMIT 1
	`

	if err := postgres.Db.Get(&post, query, postId); err != nil {
		return Post{}, err
	}

	return prepareOne(post), nil
}

func (postgres Postgres) Many(limit int, cursor uid.UID) ([]Post, error) {
	var posts []PostgresPost

	query := `
	SELECT uuid,
	       user_uuid,
	       content,
	       created_at,
	       updated_at,
	       deleted_at 
	FROM posts 
	WHERE uuid > $1 
	ORDER BY created_at 
	DESC LIMIT $2
	`

	if err := postgres.Db.Select(&posts, query, cursor, limit); err != nil {
		return nil, err
	}

	return prepareMany(posts), nil
}

func (postgres Postgres) Insert(userId uid.UID, fields PostFields) (uid.UID, error) {
	var uuid uid.UID

	query := `
	INSERT INTO posts (user_uuid, content) 
	VALUES ($1, $2) 
	RETURNING uuid
	`

	if err := postgres.Db.Get(&uuid, query, userId, fields.Content); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (postgres Postgres) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	var uuid uid.UID

	query := `
	UPDATE posts 
	SET content = $2 
	WHERE uuid = $1 
	RETURNING uuid
	`

	if err := postgres.Db.Get(&uuid, query, postId, fields.Content); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func prepareOne(pp PostgresPost) Post {
	return Post{
		Id:        pp.Uuid,
		UserId:    pp.UserUuid,
		Content:   pp.Content,
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
