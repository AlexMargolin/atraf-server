package app

import (
	"errors"
	"fmt"
	"os"
	"strings"
)

var RequiredKeys = []string{
	"MYSQL_HOST",
	"MYSQL_PORT",
	"MYSQL_USER",
	"MYSQL_PASS",
	"MYSQL_NAME",
	"SERVER_PORT",
	"SERVER_HOST",
	"ACCESS_TOKEN_SECRET",
}

// CheckEnvironment check whether required environment variables are defined.
func CheckEnvironment() error {
	missing := make([]string, 0)

	for _, key := range RequiredKeys {
		if _, exists := os.LookupEnv(key); !exists {
			missing = append(missing, key)
		}
	}

	if len(missing) != 0 {
		return errors.New(fmt.Sprintf("undefined [%s] environment variables", strings.Join(missing, ", ")))
	}

	return nil
}
