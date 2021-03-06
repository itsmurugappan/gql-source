/*
Copyright 2021 Murugappan Sevugan Chetty
*/

// Code generated by client-gen. DO NOT EDIT.

package fake

import (
	v1alpha1 "github.com/itsmurugappan/gql-source/pkg/client/clientset/versioned/typed/sources/v1alpha1"
	rest "k8s.io/client-go/rest"
	testing "k8s.io/client-go/testing"
)

type FakeSourcesV1alpha1 struct {
	*testing.Fake
}

func (c *FakeSourcesV1alpha1) GqlSources(namespace string) v1alpha1.GqlSourceInterface {
	return &FakeGqlSources{c, namespace}
}

// RESTClient returns a RESTClient that is used to communicate
// with API server by this client implementation.
func (c *FakeSourcesV1alpha1) RESTClient() rest.Interface {
	var ret *rest.RESTClient
	return ret
}
