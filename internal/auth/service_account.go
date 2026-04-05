package auth

import (
	"fmt"
	"os"

	"google.golang.org/api/option"
)

func ServiceAccountOption(credentialsPath string) (option.ClientOption, error) {
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read credentials file: %w", err)
	}
	return option.WithCredentialsJSON(data), nil
}
