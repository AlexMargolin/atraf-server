package account

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"atraf-server/pkg/uid"
)

type Account struct {
	Id             uid.UID   `json:"-"`
	Email          string    `json:"-"`
	PasswordHash   []byte    `json:"-"`
	ActivationCode string    `json:"-"`
	Active         bool      `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type Storage interface {
	NewAccount(email string, passwordHash []byte) (Account, error)
	ByEmail(email string) (Account, error)
	ByAccountId(accountId uid.UID) (Account, error)
	SetPending(accountId uid.UID) (string, error)
	SetActive(accountId uid.UID, activationCode string) error
	UpdatePassword(accountId uid.UID, passwordHash []byte) error
}

type Service struct {
	storage Storage
}

func (s Service) Register(email string, password string) (Account, error) {
	passwordHash, err := s.newPasswordHash(password)
	if err != nil {
		return Account{}, err
	}

	account, err := s.storage.NewAccount(email, passwordHash)
	if err != nil {
		return Account{}, err
	}

	if err = SendActivationMail(account.Email, account.ActivationCode); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (s Service) Login(email string, password string) (Account, error) {
	account, err := s.storage.ByEmail(email)
	if err != nil {
		return Account{}, err
	}

	if err = s.comparePasswordHash(password, account.PasswordHash); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (s Service) Forgot(email string) error {
	account, err := s.storage.ByEmail(email)
	if err != nil {
		return err
	}

	return SendPasswordResetMail(account)
}

func (s Service) Activate(accountId uid.UID, activationCode string) error {
	return s.storage.SetActive(accountId, activationCode)
}

func (s Service) Pending(accountId uid.UID) (string, error) {
	return s.storage.SetPending(accountId)
}

func (s Service) UpdatePassword(accountId uid.UID, password string) error {
	passwordHash, err := s.newPasswordHash(password)
	if err != nil {
		return err
	}

	account, err := s.storage.ByAccountId(accountId)
	if err != nil {
		return err
	}

	if err = s.storage.UpdatePassword(account.Id, passwordHash); err != nil {
		return err
	}

	if err = SendPasswordNotificationEmail(account.Email); err != nil {
		return err
	}

	return nil
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
