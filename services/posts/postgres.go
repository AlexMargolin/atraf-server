package posts

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/jmoiron/sqlx"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/uid"
	"atraf-server/services/bucket"
)

type PostgresPost struct {
	Uuid       uid.UID      `db:"uuid"`
	UserUuid   uid.UID      `db:"user_uuid"`
	Title      string       `db:"title"`
	Body       string       `db:"body"`
	Attachment string       `db:"attachment"`
	CreatedAt  time.Time    `db:"created_at"`
	UpdatedAt  sql.NullTime `db:"updated_at"`
	DeletedAt  sql.NullTime `db:"deleted_at"`
}

type Postgres struct {
	db     *sqlx.DB
	bucket *bucket.Service
}

func (p Postgres) One(postId uid.UID) (Post, error) {
	var post PostgresPost

	query := `SELECT * FROM posts WHERE posts.uuid = $1 LIMIT 1`

	// Returns an error when no results are found.
	if err := p.db.Get(&post, query, postId); err != nil {
		return Post{}, err
	}

	return p.prepareOne(post), nil
}

func (p Postgres) Many(pc *middleware.PaginationContext) ([]Post, error) {
	var posts []PostgresPost

	if pc.Cursor.Key != uid.Nil {
		query := `
		SELECT *
		FROM posts 
		WHERE (posts.created_at, posts.uuid) < ($1 :: timestamp, $2) 
		ORDER BY posts.created_at DESC 
		LIMIT $3`

		if err := p.db.Select(&posts, query, pc.Cursor.Value, pc.Cursor.Key, pc.Limit); err != nil {
			return nil, err
		}
	} else {
		query := `
		SELECT *
		FROM posts
		ORDER BY posts.created_at DESC 
		LIMIT $1`

		if err := p.db.Select(&posts, query, pc.Limit); err != nil {
			return nil, err
		}
	}

	return p.prepareMany(posts), nil
}

func (p Postgres) Insert(userId uid.UID, f *Fields) (uid.UID, error) {
	var uuid uid.UID

	attachmentPath, err := p.bucket.Save(f.File)
	if err != nil {
		return uuid, err
	}

	query := "INSERT INTO posts (user_uuid, title, body, attachment) VALUES ($1, $2, $3, $4) RETURNING uuid"
	if err = p.db.Get(&uuid, query, userId, f.Title, f.Body, attachmentPath); err != nil {
		return uuid, err
	}

	return uuid, nil
}

func (p Postgres) Update(postId uid.UID, f *Fields) error {
	query := `UPDATE posts SET title = $2, body = $3 WHERE uuid = $1`
	result, err := p.db.Exec(query, postId, f.Title, f.Body)
	if err != nil {
		return err
	}

	rows, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rows == 0 {
		return errors.New(fmt.Sprintf("no updates were made to post id [%s]", postId))
	}

	return nil
}

func (p Postgres) prepareOne(pp PostgresPost) Post {
	return Post{
		Id:         pp.Uuid,
		UserId:     pp.UserUuid,
		Title:      pp.Title,
		Body:       pp.Body,
		Attachment: p.bucket.FileURL(pp.Attachment),
		CreatedAt:  pp.CreatedAt,
		UpdatedAt:  pp.UpdatedAt.Time,
	}
}

func (p Postgres) prepareMany(pp []PostgresPost) []Post {
	var posts = make([]Post, 0)

	for _, post := range pp {
		posts = append(posts, p.prepareOne(post))
	}

	return posts
}

func NewStorage(db *sqlx.DB, b *bucket.Service) *Postgres {
	return &Postgres{db, b}
}
