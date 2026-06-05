package sdkcmd

import (
	"github.com/paymog/groundcover-cli/internal/config"
	groundcover "github.com/groundcover-com/groundcover-sdk-go"
	"github.com/groundcover-com/groundcover-sdk-go/pkg/client"
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
