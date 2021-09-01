package account

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"quotes/pkg/uid"
)

type Account struct {
	Id           uid.UID
	Email        string
	PasswordHash []byte
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Storage interface {
	ByEmail(email string) (Account, error)
	Insert(email string, passwordHash []byte) (uid.UID, error)
}

type Service struct {
	storage Storage
}

func (service *Service) Register(email string, password string) (uid.UID, error) {
	var accountId uid.UID

	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return accountId, err
	}

	accountId, err = service.storage.Insert(email, passwordHash)
	if err != nil {
		return accountId, err
	}

	return accountId, nil
}

func (service *Service) Login(email string, password string) (Account, error) {
	var account Account

	account, err := service.storage.ByEmail(email)
	if err != nil {
		return account, err
	}

	err = service.comparePasswordHash(password, account.PasswordHash)
	if err != nil {
		return account, err
	}

	return account, nil
}

func (Service) newPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (Service) comparePasswordHash(password string, passwordHash []byte) error {
	return bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
