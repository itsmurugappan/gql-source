package v1alpha1

import (
	"context"

	"knative.dev/pkg/apis"
)

// Validate implements apis.Validatable
func (g *GqlSource) Validate(ctx context.Context) *apis.FieldError {
	return g.Spec.Validate(ctx).ViaField("spec")
}

// Validate implements apis.Validatable
func (gs *GqlSourceSpec) Validate(ctx context.Context) *apis.FieldError {
	if gs.GqlServer == "" {
		return apis.ErrMissingField("GraphQL Server endpoint")
	}
	if len(gs.SubscriptionQueries) == 0 {
		return apis.ErrMissingField("No Subscription queries provided")
	}
	return nil
}
