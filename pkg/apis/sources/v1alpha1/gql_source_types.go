package v1alpha1

import (
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"knative.dev/pkg/apis"
	duckv1 "knative.dev/pkg/apis/duck/v1"
	"knative.dev/pkg/kmeta"
)

// +genclient
// +genreconciler
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GqlSource is knative eventing source for graph ql subscription
type GqlSource struct {
	metav1.TypeMeta `json:",inline"`
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// Spec holds the desired state of the GqlSource (from the client).
	// +optional
	Spec GqlSourceSpec `json:"spec,omitempty"`

	// Status communicates the observed state of the GqlSource (from the controller).
	// +optional
	Status GqlSourceStatus `json:"status,omitempty"`
}

const (
	EventType = "dev.knative.gql.event"
)

// GqlEventSource returns the GQL CloudEvent source.
func GqlEventSource(namespace, name string) string {
	return fmt.Sprintf("/apis/v1/namespaces/%s/gqlsources/%s", namespace, name)
}

var (
	// Check that GqlSource can be validated and defaulted.
	_ apis.Validatable   = (*GqlSource)(nil)
	_ apis.Defaultable   = (*GqlSource)(nil)
	_ kmeta.OwnerRefable = (*GqlSource)(nil)
	// Check that the type conforms to the duck Knative Resource shape.
	_ duckv1.KRShaped = (*GqlSource)(nil)
	_ runtime.Object  = (*GqlSource)(nil)
)

// GqlSourceSpec holds the desired state of the GqlSource (from the client).
type GqlSourceSpec struct {
	// GqlServer holds the graphql endpoint
	GqlServer string `json:"gqlServer"`
	// Subscriptions holds the list of subscription queries that need to be tracked
	SubscriptionQueries []string `json:"subscriptionQueries"`
	// Sink holds the callable where cloud events need to be sent
	Sink *duckv1.Destination `json:"sink,omitempty"`
	// GqlAuthSpec embeds the security
	GqlAuthSpec `json:",inline"`
}

type OIDCSpec struct {
	Enable bool `json:"enable,omitempty"`

	// TokenUrl is oidc token url
	// +optional
	TokenUrl SecretValueFromSource `json:"tokenUrl,omitempty"`
	// Key is the Kubernetes secret containing the client key.
	// +optional
	OIDCClient SecretValueFromSource `json:"oidcClient,omitempty"`
	// user to get oidc token.
	// +optional
	User SecretValueFromSource `json:"user,omitempty"`
	// password to get oidc token.
	// +optional
	Password SecretValueFromSource `json:"password,omitempty"`
}

// SecretValueFromSource represents the source of a secret value
type SecretValueFromSource struct {
	// The Secret key to select from.
	SecretKeyRef *corev1.SecretKeySelector `json:"secretKeyRef,omitempty"`
}

type GqlAuthSpec struct {
	OIDC OIDCSpec `json:"oidc,omitempty"`
}

type GqlSourceStatus struct {
	// inherits duck/v1 SourceStatus, which currently provides:
	// * ObservedGeneration - the 'Generation' of the Service that was last
	//   processed by the controller.
	// * Conditions - the latest available observations of a resource's current
	//   state.
	// * SinkURI - the current active sink URI that has been configured for the
	//   Source.
	duckv1.SourceStatus `json:",inline"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// GqlSourceList is a list of GqlSource resources
type GqlSourceList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata"`

	Items []GqlSource `json:"items"`
}

// GetStatus retrieves the status of the resource. Implements the KRShaped interface.
func (as *GqlSource) GetStatus() *duckv1.Status {
	return &as.Status.Status
}
