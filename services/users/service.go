package users

import (
	"time"

	"atraf-server/pkg/uid"
)

type User struct {
	Id             uid.UID   `json:"id"`
	Email          string    `json:"email"`
	FirstName      string    `json:"first_name"`
	LastName       string    `json:"last_name"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"created_at"`
}

// UserFields is a struct representing all Post values
// which can be modified by the client.
type UserFields struct {
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	ProfilePicture string `json:"profile_picture"`
	Email          string `json:"email" validate:"required,email"`
}

type Storage interface {
	ById(userId uid.UID) (User, error)
	ByIds(userIds []uid.UID) ([]User, error)
	ByAccountId(accountID uid.UID) (User, error)
	Insert(accountId uid.UID, fields UserFields) error
}

type Service struct {
	storage Storage
}

func (s Service) UserById(userId uid.UID) (User, error) {
	return s.storage.ById(userId)
}

func (s Service) UsersByIds(userIds []uid.UID) ([]User, error) {
	return s.storage.ByIds(userIds)
}

func (s Service) UserByAccountId(accountId uid.UID) (User, error) {
	return s.storage.ByAccountId(accountId)
}

func (s Service) NewUser(accountId uid.UID, fields UserFields) error {
	return s.storage.Insert(accountId, fields)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
