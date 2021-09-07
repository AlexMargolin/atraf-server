package users

import (
	"atraf-server/pkg/uid"
)

type User struct {
	Id             uid.UID `json:"id"`
	Email          string  `json:"email"`
	FirstName      string  `json:"first_name"`
	LastName       string  `json:"last_name"`
	ProfilePicture string  `json:"profile_picture"`
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
	Insert(accountId uid.UID, fields UserFields) (uid.UID, error)
}

type Service struct {
	storage Storage
}

func (service *Service) UserById(userId uid.UID) (User, error) {
	return service.storage.ById(userId)
}

func (service *Service) UsersByIds(userIds []uid.UID) ([]User, error) {
	return service.storage.ByIds(userIds)
}

func (service *Service) UserByAccount(accountId uid.UID) (User, error) {
	return service.storage.ByAccountId(accountId)
}

func (service *Service) NewUser(accountId uid.UID, fields UserFields) (uid.UID, error) {
	return service.storage.Insert(accountId, fields)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
