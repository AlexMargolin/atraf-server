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

// Fields is a struct representing all Post values
// which can be modified by the client.
type Fields struct {
	Title string         `json:"title" validate:"required"`
	Body  string         `json:"body" validate:"required"`
	File  multipart.File `json:"file" validate:"required"`
}

type Storage interface {
	One(postId uid.UID) (Post, error)
	Many(pagination *middleware.PaginationContext) ([]Post, error)
	Insert(userId uid.UID, fields *Fields) (uid.UID, error)
	Update(postId uid.UID, fields *Fields) error
}

type Service struct {
	storage Storage
}

func (s Service) NewPost(userId uid.UID, f *Fields) (uid.UID, error) {
	return s.storage.Insert(userId, f)
}

func (s Service) PostById(postId uid.UID) (Post, error) {
	return s.storage.One(postId)
}

func (s Service) Posts(p *middleware.PaginationContext) ([]Post, error) {
	return s.storage.Many(p)
}

func (s Service) UpdatePost(postId uid.UID, f *Fields) error {
	return s.storage.Update(postId, f)
}

func UniqueUserIds(p []Post) []uid.UID {
	userIds := make([]uid.UID, 0)
	m := make(map[uid.UID]bool, 0)

	for _, post := range p {
		if m[post.UserId] {
			continue
		}
		m[post.UserId] = true
		userIds = append(userIds, post.UserId)
	}

	return userIds
}

func NewService(s Storage) *Service {
	return &Service{s}
}
