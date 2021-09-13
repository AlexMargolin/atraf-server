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

func (service *Service) Register(email string, password string) (Account, error) {
	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return Account{}, err
	}

	account, err := service.storage.NewAccount(email, passwordHash)
	if err != nil {
		return Account{}, err
	}

	if err = SendActivationMail(account.Email, account.ActivationCode); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (service *Service) Login(email string, password string) (Account, error) {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return Account{}, err
	}

	if err = service.comparePasswordHash(password, account.PasswordHash); err != nil {
		return Account{}, err
	}

	return account, nil
}

func (service *Service) Forgot(email string) error {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return err
	}

	return SendPasswordResetMail(account)
}

func (service *Service) Activate(accountId uid.UID, activationCode string) error {
	return service.storage.SetActive(accountId, activationCode)
}

func (service *Service) Pending(accountId uid.UID) (string, error) {
	return service.storage.SetPending(accountId)
}

func (service *Service) UpdatePassword(accountId uid.UID, password string) error {
	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return err
	}

	account, err := service.storage.ByAccountId(accountId)
	if err != nil {
		return err
	}

	if err = service.storage.UpdatePassword(account.Id, passwordHash); err != nil {
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
