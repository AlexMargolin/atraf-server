package account

import (
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
)

const (
	AccessTokenValidFor = time.Minute * 15
	ResetTokenValidFor  = time.Minute * 5
)

type Account struct {
	Id           uid.UID   `json:"-"`
	Email        string    `json:"-"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Storage interface {
	ById(accountId uid.UID) (Account, error)
	ByEmail(email string) (Account, error)
	Insert(email string, passwordHash []byte) (uid.UID, int, error)
	UpdatePassword(accountId uid.UID, passwordHash []byte) error
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

	accountId, code, err := service.storage.Insert(email, passwordHash)
	if err != nil {
		return accountId, err
	}

	if err = ActivationEmail(email, code); err != nil {
		return accountId, err
	}

	return accountId, nil
}

func (service *Service) Login(email string, password string) (string, error) {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return "", err
	}

	if err = service.comparePasswordHash(password, account.PasswordHash); err != nil {
		return "", err
	}

	claims := token.Claims{
		Subject:   account.Id.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(AccessTokenValidFor).Unix(),
	}

	secret := os.Getenv("ACCESS_TOKEN_SECRET")
	accessToken, err := token.New(secret, claims)
	if err != nil {
		return "", err
	}

	return accessToken, nil
}

func (service *Service) Forgot(email string) error {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return err
	}

	claims := token.Claims{
		Subject:   account.Id.String(),
		IssuedAt:  time.Now().Unix(),
		ExpiresAt: time.Now().Add(ResetTokenValidFor).Unix(),
	}

	secret := os.Getenv("RESET_TOKEN_SECRET")
	resetToken, err := token.New(secret, claims)
	if err != nil {
		return err
	}

	if err = PasswordResetMail(account.Email, resetToken); err != nil {
		return err
	}

	return nil
}

func (service *Service) UpdatePassword(accountId uid.UID, password string) error {
	account, err := service.storage.ById(accountId)
	if err != nil {
		return err
	}

	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return err
	}

	if err = service.storage.UpdatePassword(accountId, passwordHash); err != nil {
		return err
	}

	if err = PasswordNotification(account.Email); err != nil {
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
