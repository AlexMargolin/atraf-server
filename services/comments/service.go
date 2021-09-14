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
	Insert(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields CommentFields) (Comment, error)
	Update(commentId uid.UID, fields CommentFields) error
}

type Service struct {
	storage Storage
}

func (service *Service) NewComment(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields CommentFields) (Comment, error) {
	return service.storage.Insert(userId, sourceId, parentId, fields)
}

func (service *Service) UpdateComment(commentId uid.UID, fields CommentFields) error {
	return service.storage.Update(commentId, fields)
}

func (service *Service) CommentsBySourceId(sourceId uid.UID) ([]Comment, error) {
	return service.storage.Many(sourceId)
}

func UniqueUserIds(comments []Comment) []uid.UID {
	userIds := make([]uid.UID, 0)
	m := make(map[uid.UID]bool, 0)

	for _, comment := range comments {
		if m[comment.UserId] {
			continue
		}
		m[comment.UserId] = true
		userIds = append(userIds, comment.UserId)
	}

	return userIds
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
