package account

import (
	"time"

	"golang.org/x/crypto/bcrypt"

	"quotes/pkg/uid"
)

type Account struct {
	Id        uid.UID
	Email     string
	Password  []byte
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}

type Storage interface {
	ByEmail(email string) (*Account, error)
	Insert(email string, password []byte) (uid.UID, error)
}

type Service struct {
	storage Storage
}

// Register creates a new account using the provided email and password arguments
func (service *Service) Register(email string, password string) (uid.UID, error) {
	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return uid.Nil, err
	}

	accountId, err := service.storage.Insert(email, passwordHash)
	if err != nil {
		return uid.Nil, err
	}

	return accountId, nil
}

// Login attempts to fetch an existing account by the email address.
// if found, the password argument is then compared to the Accounts stored hash value.
// returns an Account on successful verification.
func (service *Service) Login(email string, password string) (*Account, error) {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return nil, err
	}

	err = service.comparePasswordHash(password, account.Password)
	if err != nil {
		return nil, err
	}

	return account, nil
}

// newPasswordHash receives a plain text password and returns a bcrypt password hash.
func (Service) newPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

// comparePasswordHash receives a plain text password, and a bcrypt password hash.
// Returns nil on successful compare, error otherwise
func (Service) comparePasswordHash(password string, passwordHash []byte) error {
	return bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
}

// NewService returns a new Account Service instance
func NewService(storage Storage) *Service {
	return &Service{storage}
}
