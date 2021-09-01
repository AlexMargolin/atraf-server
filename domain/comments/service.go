package comments

import (
	"time"

	"quotes/pkg/uid"
)

type Comment struct {
	Id        uid.UID   `json:"id"`
	UserId    uid.UID   `json:"user_id"`
	PostId    uid.UID   `json:"post_id"`
	ParentId  uid.UID   `json:"parent_id"`
	Content   string    `json:"content"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentFields is a struct representing all Comment values
// which can be modified by the client.
type CommentFields struct {
	Content string `json:"content" validate:"required"`
}

type Storage interface {
	Many(postId uid.UID) ([]Comment, error)
	Insert(userId uid.UID, postId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error)
	Update(commentId uid.UID, fields CommentFields) (uid.UID, error)
}

type Service struct {
	storage Storage
}

func (service *Service) New(userId uid.UID, postId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	return service.storage.Insert(userId, postId, parentId, fields)
}

func (service *Service) Update(commentId uid.UID, fields CommentFields) (uid.UID, error) {
	return service.storage.Update(commentId, fields)
}

func (service *Service) Comments(postId uid.UID) ([]Comment, error) {
	return service.storage.Many(postId)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
