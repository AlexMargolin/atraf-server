package posts

import (
	"time"

	"quotes/pkg/uid"
)

type Post struct {
	Id        uid.UID   `json:"id"`
	UserId    uid.UID   `json:"user_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// PostFields is a struct representing all Post values
// which can be modified by the client.
type PostFields struct {
	Content string `json:"content" validate:"required"`
}

type Storage interface {
	One(postId uid.UID) (Post, error)
	Many(limit int, cursor uid.UID) ([]Post, error)
	Update(postId uid.UID, fields PostFields) (uid.UID, error)
	Insert(userId uid.UID, fields PostFields) (uid.UID, error)
}

type Service struct {
	storage Storage
}

func (service *Service) Post(postId uid.UID) (Post, error) {
	return service.storage.One(postId)
}

func (service *Service) Posts(limit int, cursor uid.UID) ([]Post, error) {
	return service.storage.Many(limit, cursor)
}

func (service *Service) New(userId uid.UID, fields PostFields) (uid.UID, error) {
	return service.storage.Insert(userId, fields)
}

func (service *Service) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	return service.storage.Update(postId, fields)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
