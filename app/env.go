package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

// CheckEnvironment ensures the required environment variables are defined and not empty.
func CheckEnvironment() error {
	missing := make([]string, 0)

	var RequiredKeys = []string{
		"CLIENT_URL",
		"DB_HOST",
		"DB_PORT",
		"DB_USER",
		"DB_PASS",
		"DB_NAME",
		"SERVER_PORT",
		// "SERVER_HOST",
		"SMTP_HOST",
		"SMTP_PORT",
		"SMTP_USER",
		"SMTP_PASS",
		"BUCKET_URL",
		"ACCESS_TOKEN_SECRET",
		"RESET_TOKEN_SECRET",
	}

	for _, key := range RequiredKeys {
		if value := os.Getenv(key); value == "" {
			missing = append(missing, key)
		}
	}

	if len(missing) != 0 {
		return errors.New(fmt.Sprintf("undefined [%s] environment variables", strings.Join(missing, ", ")))
	}

	return nil
}
