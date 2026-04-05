package auth

import "google.golang.org/api/option"

func ServiceAccountOption(credentialsPath string) (option.ClientOption, error) {
	return option.WithCredentialsFile(credentialsPath), nil
}
