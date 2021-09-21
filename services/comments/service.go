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

// Fields is a struct representing all Comment values
// which can be modified by the client.
type Fields struct {
	Body string `json:"body" validate:"required"`
}

type Storage interface {
	Insert(userId uid.UID, sourceId uid.UID, parentId uid.UID, data *Fields) (Comment, error)
	Update(commentId uid.UID, data *Fields) error
	Many(sourceId uid.UID) ([]Comment, error)
}

type Service struct {
	storage Storage
}

func (s Service) NewComment(userId uid.UID, sourceId uid.UID, parentId uid.UID, fields *Fields) (Comment, error) {
	return s.storage.Insert(userId, sourceId, parentId, fields)
}

func (s Service) UpdateComment(commentId uid.UID, fields *Fields) error {
	return s.storage.Update(commentId, fields)
}

func (s Service) CommentsBySourceId(sourceId uid.UID) ([]Comment, error) {
	return s.storage.Many(sourceId)
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
