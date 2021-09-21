package account

import (
	"fmt"
	"net/mail"
	"os"
	"time"

	"golang.org/x/crypto/bcrypt"

	"atraf-server/pkg/mailer"
	"atraf-server/pkg/token"
	"atraf-server/pkg/uid"
)

type Account struct {
	Id             uid.UID   `json:"-"`
	Email          string    `json:"-"`
	PasswordHash   []byte    `json:"-"`
	ActivationCode string    `json:"-"`
	Active         bool      `json:"active"`
	Nickname       string    `json:"-"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

type Storage interface {
	Insert(email string, nickname string, passwordHash []byte) (Account, error)
	ByEmail(email string) (Account, error)
	ByAccountId(accountId uid.UID) (Account, error)
	SetPending(accountId uid.UID) (string, error)
	SetActive(accountId uid.UID, activationCode string) error
	UpdatePassword(accountId uid.UID, passwordHash []byte) error
}

type Service struct {
	storage Storage
}

func (s Service) ByAccountId(accountId uid.UID) (Account, error) {
	return s.storage.ByAccountId(accountId)
}

func (s Service) Register(password string, email string, nickname string) (Account, error) {
	passwordHash, err := s.newPasswordHash(password)
	if err != nil {
		return Account{}, err
	}

	account, err := s.storage.Insert(email, nickname, passwordHash)
	if err != nil {
		return Account{}, err
	}

	if err = s.sendActivationMail(account); err != nil {
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

	return s.sendPasswordResetMail(account)
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

	if err = s.passwordNotificationMail(account); err != nil {
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

func (Service) sendPasswordResetMail(account Account) error {
	subject := "Password Reset Request"
	from := mail.Address{
		Name:    "Atraf Accounts",
		Address: "accounts@atraf.app",
	}

	resetToken, err := token.NewResetToken(token.ResetTokensCustomClaims{
		AccountId: account.Id,
	})
	if err != nil {
		return err
	}

	data := struct {
		ResetURL string
		Duration float64
	}{
		ResetURL: fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_URL"), resetToken),
		Duration: token.ResetTokenValidFor.Minutes(),
	}

	filename := "templates/password-reset.html"
	return mailer.FromTemplate(filename, data, subject, from, []string{account.Email})
}

func (Service) sendActivationMail(account Account) error {
	subject := "Account Activation Code"
	from := mail.Address{
		Name:    "Atraf Accounts",
		Address: "accounts@atraf.app",
	}

	data := struct {
		Code string
	}{
		Code: account.ActivationCode,
	}

	filename := "templates/account-activation.html"
	return mailer.FromTemplate(filename, data, subject, from, []string{account.Email})
}

func (Service) passwordNotificationMail(account Account) error {
	subject := "Password reset notification"
	from := mail.Address{
		Name:    "Atraf Accounts",
		Address: "accounts@atraf.app",
	}

	filename := "templates/password-change-notice.html"
	return mailer.FromTemplate(filename, nil, subject, from, []string{account.Email})
}

func NewService(storage Storage) *Service {
	return &Service{storage}
}
