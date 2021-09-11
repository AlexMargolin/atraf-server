package account

import (
	"fmt"
	"net/mail"
	"os"

	"atraf-server/pkg/mailer"
)

const (
	ActivationTemplate     = "templates/account-activation.html"
	PasswordChangeTemplate = "templates/password-change-notice.html"
	PasswordResetTemplate  = "templates/password-reset.html"
)

func ActivationEmail(to string, code int) error {
	type ActivationData struct {
		Code int
	}

	subject := "Account Activation Code"

	from := mail.Address{
		Name:    "Atraf Support",
		Address: "support@atraf.app",
	}

	data := &ActivationData{
		code,
	}

	return mailer.FromTemplate(ActivationTemplate, data, subject, from, []string{to})
}

func PasswordResetMail(to string, resetToken string) error {
	type PasswordResetData struct {
		ResetURL string
		Duration float64
	}

	subject := "Password Reset Request"

	from := mail.Address{
		Name:    "Atraf Support",
		Address: "support@atraf.app",
	}

	url := fmt.Sprintf("%s/reset/%s", os.Getenv("CLIENT_URL"), resetToken)
	data := &PasswordResetData{
		url,
		ResetTokenValidFor.Minutes(),
	}

	return mailer.FromTemplate(PasswordResetTemplate, data, subject, from, []string{to})
}

func PasswordNotification(to string) error {
	subject := "Account password was reset"

	from := mail.Address{
		Name:    "Atraf Support",
		Address: "support@atraf.app",
	}

	return mailer.FromTemplate(PasswordChangeTemplate, nil, subject, from, []string{to})
}
