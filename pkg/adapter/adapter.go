package adapter

import (
	"context"

	"go.uber.org/zap"
	"knative.dev/eventing/pkg/adapter/v2"
	"knative.dev/pkg/logging"

	cloudevents "github.com/cloudevents/sdk-go/v2"
)

type adapterConfig struct {
	adapter.EnvConfig

	GqlServer  string   `envconfig:"GQL_SERVER" required:"true"`
	SubQueries []string `envconfig:"SUBSCRIPTION_QUERIES" required:"true"`
}

func NewEnvConfig() adapter.EnvConfigAccessor {
	return &adapterConfig{}
}

type Adapter struct {
	config *adapterConfig
	logger *zap.Logger
	client cloudevents.Client
}

func NewAdapter(ctx context.Context, processed adapter.EnvConfigAccessor, ceClient cloudevents.Client) adapter.Adapter {
	logger := logging.FromContext(ctx).Desugar()
	config := processed.(*adapterConfig)

	return &Adapter{
		config: config,
		client: ceClient,
		logger: logger,
	}
}

func (a *Adapter) Start(ctx context.Context) error {
	return a.start(ctx.Done())
}

func (a *Adapter) start(stopCh <-chan struct{}) error {
	ctx, cancel := context.WithCancel(context.Background())

	//subscribe to each query in separate go routines
	for _, q := range a.config.SubQueries {
		go a.subscribe(ctx, q)
	}

	<-stopCh
	cancel()
	a.logger.Info("Shutting down...")
	return nil
}
