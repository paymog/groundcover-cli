package sdkcmd

import (
	"context"
	"fmt"

	"github.com/paymog/groundcover-cli/internal/config"
	groundcover "github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client/apikeys"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/option"
)

func newClient(cfg config.Config) (*client.GroundcoverAPI, error) {
	cfg.ApplyEnv()
	if err := cfg.RequireSDKAuth(); err != nil {
		return nil, err
	}
	return groundcover.NewClient(
		option.WithAPIKey(cfg.APIKey),
		option.WithBackendID(cfg.BackendID),
		option.WithBaseURL(cfg.NormalizedBaseURL()),
	)
}

// ValidateAuth confirms that the API key and backend ID in cfg can make an
// authenticated request, by issuing a cheap read (listing API keys). It is used
// by `auth login` to verify credentials before persisting them.
func ValidateAuth(ctx context.Context, cfg config.Config) error {
	c, err := newClient(cfg)
	if err != nil {
		return err
	}
	if _, err := c.Apikeys.ListAPIKeys(apikeys.NewListAPIKeysParams().WithContext(ctx), nil); err != nil {
		return fmt.Errorf("credential check failed (verify the API key and backend ID): %w", err)
	}
	return nil
}
