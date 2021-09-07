package comments

import (
	"time"

	"atraf-server/pkg/uid"
)

type Comment struct {
	Id        uid.UID   `json:"id"`
	UserId    uid.UID   `json:"user_id"`
	SourceId  uid.UID   `json:"source_id"`
	ParentId  uid.UID   `json:"parent_id"`
	Body      string    `json:"body"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// CommentFields is a struct representing all Comment values
// which can be modified by the client.
type CommentFields struct {
	Body string `json:"body" validate:"required"`
}

type Storage interface {
	Many(sourceId uid.UID) ([]Comment, error)
	Insert(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error)
	Update(commentId uid.UID, fields CommentFields) (uid.UID, error)
}

type Service struct {
	storage Storage
}

func (service *Service) New(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields CommentFields) (uid.UID, error) {
	return service.storage.Insert(userId, sourceId, parentId, fields)
}

func (service *Service) Update(commentId uid.UID, fields CommentFields) (uid.UID, error) {
	return service.storage.Update(commentId, fields)
}

func (service *Service) Comments(sourceId uid.UID) ([]Comment, error) {
	return service.storage.Many(sourceId)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
