package account

import (
	"fmt"
	"log"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
)

type Account struct {
	Id           uid.UID   `json:"-"`
	Email        string    `json:"-"`
	PasswordHash []byte    `json:"-"`
	CreatedAt    time.Time `json:"-"`
	UpdatedAt    time.Time `json:"-"`
}

type Storage interface {
	ByEmail(email string) (Account, error)
	Insert(email string, passwordHash []byte) (uid.UID, error)
	SetReset(accountId uid.UID) (uid.UID, error)
	UpdatePassword(tokenId uid.UID, passwordHash []byte) error
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

func (service *Service) Login(email string, password string) (string, error) {
	account, err := service.storage.ByEmail(email)
	if err != nil {
		return "", err
	}

	if err = service.comparePasswordHash(password, account.PasswordHash); err != nil {
		return "", err
	}

	accessToken, err := token.New(os.Getenv("ACCESS_TOKEN_SECRET"), account.Id.String())
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

	tokenId, err := service.storage.SetReset(account.Id)
	if err != nil {
		return err
	}

	resetToken, err := token.New(os.Getenv("RESET_TOKEN_SECRET"), tokenId.String())
	if err != nil {
		return err
	}

	resetLink := service.resetLink(resetToken)
	log.Println(resetLink) // todo email the link

	return nil
}

func (service *Service) Reset(unverifiedToken string, password string) error {
	passwordHash, err := service.newPasswordHash(password)
	if err != nil {
		return err
	}

	claims, err := token.Verify(os.Getenv("RESET_TOKEN_SECRET"), unverifiedToken)
	if err != nil {
		return err
	}

	tokenId, err := uid.FromString(claims.Subject)
	if err != nil {
		return err
	}

	if err = service.storage.UpdatePassword(tokenId, passwordHash); err != nil {
		return err
	}

	// todo send password changed notification email
	log.Println("password changed")
	return nil
}

func (Service) newPasswordHash(password string) ([]byte, error) {
	return bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
}

func (Service) comparePasswordHash(password string, passwordHash []byte) error {
	return bcrypt.CompareHashAndPassword(passwordHash, []byte(password))
}

func (Service) resetLink(resetToken string) string {
	return fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_URL"), resetToken)
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
