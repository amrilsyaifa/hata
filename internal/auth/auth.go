package auth

import (
	"context"
	"fmt"

	"github.com/amrilsyaifa/hata/internal/config"
	"google.golang.org/api/option"
)

func ClientOption(ctx context.Context, cfg *config.Config) (option.ClientOption, error) {
	switch cfg.Auth.Type {
	case "service_account":
		return ServiceAccountOption(cfg.Auth.CredentialsPath)
	case "oauth":
		return OAuthOption(ctx, cfg.Auth.CredentialsPath, cfg.Auth.TokenPath)
	default:
		return nil, fmt.Errorf("unknown auth type: %q (must be 'service_account' or 'oauth')", cfg.Auth.Type)
	}
}
