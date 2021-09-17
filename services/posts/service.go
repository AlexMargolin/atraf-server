package posts

import (
	"mime/multipart"
	"time"

	"atraf-server/pkg/middleware"
	"atraf-server/pkg/uid"
)

type Post struct {
	Id         uid.UID   `json:"id"`
	UserId     uid.UID   `json:"user_id"`
	Title      string    `json:"title"`
	Body       string    `json:"body"`
	Attachment string    `json:"attachment"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

// PostFields is a struct representing all Post values
// which can be modified by the client.
type PostFields struct {
	Title string         `validate:"required"`
	Body  string         `validate:"required"`
	File  multipart.File `validate:"required"`
}

type Storage interface {
	One(postId uid.UID) (Post, error)
	Many(pagination *middleware.PaginationContext) ([]Post, error)
	Insert(userId uid.UID, fields *PostFields) (uid.UID, error)
	Update(postId uid.UID, fields *PostFields) error
}

type Service struct {
	storage Storage
}

func (s Service) PostById(postId uid.UID) (Post, error) {
	return s.storage.One(postId)
}

func (s Service) ListPosts(pagination *middleware.PaginationContext) ([]Post, error) {
	return s.storage.Many(pagination)
}

func (s Service) NewPost(userId uid.UID, fields *PostFields) (uid.UID, error) {
	return s.storage.Insert(userId, fields)
}

func (s Service) UpdatePost(postId uid.UID, fields *PostFields) error {
	return s.storage.Update(postId, fields)
}

func UniqueUserIds(posts []Post) []uid.UID {
	userIds := make([]uid.UID, 0)
	m := make(map[uid.UID]bool, 0)

	for _, post := range posts {
		if m[post.UserId] {
			continue
		}
		m[post.UserId] = true
		userIds = append(userIds, post.UserId)
	}

	return userIds
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
