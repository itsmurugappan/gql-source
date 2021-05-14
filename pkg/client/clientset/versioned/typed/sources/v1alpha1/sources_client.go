/*
Copyright 2021 Murugappan Sevugan Chetty
*/

// Code generated by client-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
	"github.com/itsmurugappan/gql-source/pkg/client/clientset/versioned/scheme"
	rest "k8s.io/client-go/rest"
)

type SourcesV1alpha1Interface interface {
	RESTClient() rest.Interface
	GqlSourcesGetter
}

// SourcesV1alpha1Client is used to interact with features provided by the sources.muru.dev group.
type SourcesV1alpha1Client struct {
	restClient rest.Interface
}

func (c *SourcesV1alpha1Client) GqlSources(namespace string) GqlSourceInterface {
	return newGqlSources(c, namespace)
}

// NewForConfig creates a new SourcesV1alpha1Client for the given config.
func NewForConfig(c *rest.Config) (*SourcesV1alpha1Client, error) {
	config := *c
	if err := setConfigDefaults(&config); err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(&config)
	if err != nil {
		return nil, err
	}
	return &SourcesV1alpha1Client{client}, nil
}

// NewForConfigOrDie creates a new SourcesV1alpha1Client for the given config and
// panics if there is an error in the config.
func NewForConfigOrDie(c *rest.Config) *SourcesV1alpha1Client {
	client, err := NewForConfig(c)
	if err != nil {
		panic(err)
	}
	return client
}

// New creates a new SourcesV1alpha1Client for the given RESTClient.
func New(c rest.Interface) *SourcesV1alpha1Client {
	return &SourcesV1alpha1Client{c}
}

func setConfigDefaults(config *rest.Config) error {
	gv := v1alpha1.SchemeGroupVersion
	config.GroupVersion = &gv
	config.APIPath = "/apis"
	config.NegotiatedSerializer = scheme.Codecs.WithoutConversion()

	if config.UserAgent == "" {
		config.UserAgent = rest.DefaultKubernetesUserAgent()
	}

	return nil
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *SourcesV1alpha1Client) RESTClient() rest.Interface {
	if c == nil {
		return nil
	}
	return c.restClient
}
