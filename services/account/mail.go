package account

import (
	"fmt"
	"net/mail"
	"os"

	"atraf-server/pkg/mailer"
	"atraf-server/pkg/token"
)

const (
	ActivationTemplate     = "templates/account-activation.html"
	PasswordChangeTemplate = "templates/password-change-notice.html"
	PasswordResetTemplate  = "templates/password-reset.html"
)

var from = mail.Address{
	Name:    "Atraf Accounts",
	Address: "support@atraf.app",
}

func SendActivationMail(account Account) error {
	type Data struct {
		ActivationURL string
		Duration      float64
	}

	activationToken, err := token.NewActivationToken(token.ActivationTokensCustomClaims{
		AccountId: account.Id,
	})

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_URL"), activationToken)
	data := &Data{
		url,
		token.ActivationTokenValidFor.Minutes(),
	}

	return mailer.FromTemplate(ActivationTemplate, data, "New Account Activation", from, []string{account.Email})
}

func SendPasswordResetMail(account Account) error {
	type Data struct {
		ResetURL string
		Duration float64
	}

	resetToken, err := token.NewResetToken(token.ResetTokensCustomClaims{
		AccountId: account.Id,
	})

	if err != nil {
		return err
	}

	url := fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_URL"), resetToken)
	data := &Data{
		url,
		token.ResetTokenValidFor.Minutes(),
	}

	return mailer.FromTemplate(PasswordResetTemplate, data, "Password Reset Request", from, []string{account.Email})
}

func SendPasswordNotificationEmail(to string) error {
	return mailer.FromTemplate(PasswordChangeTemplate, nil, "Account password was reset", from, []string{to})
}
