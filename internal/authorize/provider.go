package authorize

import (
	context "context"
	"fmt"
	"github.com/absurdlab/tigerd/internal/healthprobe"
	providerv1 "github.com/absurdlab/tigerd/proto/gen/go/proto/provider/v1"
	"github.com/absurdlab/tigerd/proto/gen/go/proto/provider/v1/providerv1connect"
	"github.com/bufbuild/connect-go"
	validation "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
	"github.com/hellofresh/health-go/v5"
	"net/http"
	"time"
)

// ProviderProperties is the configuration properties for a single provider.
type ProviderProperties struct {
	// Key is the unique identifier of this provider.
	Key string `json:"key" yaml:"key"`
	// Address is the dial string address to connect to the provider. Non-localhost is supported, however, will
	// print a WARNING message to console.
	Address string `json:"address" yaml:"address"`
}

// Validate performs validation to this ProviderProperties.
func (p *ProviderProperties) Validate() error {
	err := validation.Errors{
		"key": validation.Validate(p.Key, validation.Required),
		"address": validation.Validate(p.Address,
			validation.Required,
			is.DialString,
		),
	}.Filter()
	if err != nil {
		return err
	}
	return nil
}

func NewProviderHealthProbes(configs []*ProviderProperties) healthprobe.Interface {
	probe := &providerHealthProbes{services: map[string]providerv1connect.PingServiceClient{}}
	for _, c := range configs {
		probe.services[c.Key] = providerv1connect.NewPingServiceClient(http.DefaultClient, c.Address)
	}
	return probe
}

type providerHealthProbes struct {
	services map[string]providerv1connect.PingServiceClient
}

func (p *providerHealthProbes) Register(h *health.Health) error {
	for key, ping := range p.services {
		if err := h.Register(health.Config{
			Name:      fmt.Sprintf("provider:%s", key),
			Timeout:   10 * time.Second,
			SkipOnErr: true,
			Check: func(ctx context.Context) error {
				_, err := ping.Ping(ctx, &connect.Request[providerv1.PingRequest]{
					Msg: &providerv1.PingRequest{},
				})
				return err
			},
		}); err != nil {
			return err
		}
	}

	return nil
}
