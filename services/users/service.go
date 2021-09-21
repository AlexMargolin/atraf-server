package users

import (
	"time"

	"atraf-server/pkg/uid"
)

type User struct {
	Id             uid.UID   `json:"id"`
	Email          string    `json:"-"`
	Nickname       string    `json:"nickname"`
	ProfilePicture string    `json:"profile_picture"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

// Fields are User fields which can be modified.
type Fields struct {
	Email          string `json:"email" validate:"required,email"`
	Nickname       string `json:"nickname" validate:"required"`
	ProfilePicture string `json:"profile_picture"`
}

type Storage interface {
	ById(userId uid.UID) (User, error)
	ByIds(userIds []uid.UID) ([]User, error)
	ByAccountId(accountID uid.UID) (User, error)
	Insert(accountId uid.UID, fields *Fields) error
}

type Service struct {
	storage Storage
}

func (s Service) NewUser(accountId uid.UID, f *Fields) error {
	return s.storage.Insert(accountId, f)
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

func NewService(storage Storage) *Service {
	return &Service{storage}
}
