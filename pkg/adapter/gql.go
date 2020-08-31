package adapter

import (
	"context"
	"encoding/json"
	"net/http"
	"time"

	"go.uber.org/zap"

	"github.com/google/uuid"

	cloudevents "github.com/cloudevents/sdk-go/v2"
	"github.com/itsmurugappan/graphql"

	sourcesv1alpha1 "github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
)

const (
	interval time.Duration = 10 * time.Second
	retry    int           = 3
)

func (a *Adapter) subscribe(ctx context.Context, q string) {
	ctx = cloudevents.ContextWithRetriesExponentialBackoff(ctx, interval, retry)

	client := graphql.NewClient(a.config.GqlServer)
	subClient, err := client.SubscriptionClient(ctx, http.Header{"Cache-Control": []string{"no-cache"}})
	if err != nil {
		a.logger.Error("Error making subscription client", zap.Error(err))
		return
	}

	req := graphql.NewRequest(q)
	sub, err := subClient.Subscribe(req)
	if err != nil {
		a.logger.Error("Error subscribing with query", zap.String("gql server", a.config.GqlServer), zap.String("query", q), zap.Error(err))
		return
	}

	for {
		select {
		case res := <-sub:
			a.postToSink(ctx, *res.Data)
		case <-ctx.Done():
			subClient.Unsubscribe(sub)
			return
		}
	}
}

func (a *Adapter) postToSink(ctx context.Context, data json.RawMessage) {
	event := cloudevents.NewEvent(cloudevents.VersionV1)

	id, _ := uuid.NewUUID()
	event.SetID(id.String())
	event.SetTime(time.Now())
	event.SetType(sourcesv1alpha1.EventType)
	event.SetSource(sourcesv1alpha1.GqlEventSource(a.config.Namespace, a.config.Name))

	if err := event.SetData(cloudevents.ApplicationJSON, message(data)); err != nil {
		a.logger.Error("Error setting data ", zap.Error(err))
		return
	}

	if err := a.client.Send(ctx, event); !cloudevents.IsACK(err) {
		a.logger.Error("failed to send gql subscription cloudevent", zap.Error(err))
	}
}

type Message struct {
	Body string `json:"body"`
}

func message(data json.RawMessage) interface{} {
	// try to marshal the body into an interface.
	var obj map[string]*json.RawMessage
	if err := json.Unmarshal(data, &obj); err != nil {
		//default to a wrapped message.
		return Message{Body: string(data)}
	}
	return obj
}
