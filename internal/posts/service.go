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
	Content string `validate:"required"`
}

type Storage interface {
	Count() (int, error)
	One(postId uid.UID) (Post, error)
	Many(offset int, limit int) ([]Post, error)
	Update(postId uid.UID, fields PostFields) (uid.UID, error)
	Insert(userId uid.UID, fields PostFields) (uid.UID, error)
}

type Service struct {
	storage Storage
}

// Post retrieves a single Post from the storage
func (service *Service) Post(postId uid.UID) (Post, error) {
	return service.storage.One(postId)
}

// Posts returns a paginated list of posts
// We Assume the page value at this point is 1 or above
func (service *Service) Posts(pageNum int, perPage int) ([]Post, error) {
	// DB Offset starts at 0
	offset := (pageNum - 1) * perPage

	return service.storage.Many(offset, perPage)
}

// New inserts a new Post to the storage
func (service *Service) New(userId uid.UID, fields PostFields) (uid.UID, error) {
	return service.storage.Insert(userId, fields)
}

// Update updates an existing post in the storage
func (service *Service) Update(postId uid.UID, fields PostFields) (uid.UID, error) {
	return service.storage.Update(postId, fields)
}

// Total counts the number of Posts in the storage
func (service *Service) Total() (int, error) {
	return service.storage.Count()
}

// NewService returns a new Posts service instance
func NewService(storage Storage) *Service {
	return &Service{storage}
}
