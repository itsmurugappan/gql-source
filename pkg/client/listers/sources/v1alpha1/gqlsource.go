/*
Copyright 2021 Murugappan Sevugan Chetty
*/

// Code generated by lister-gen. DO NOT EDIT.

package v1alpha1

import (
	v1alpha1 "github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/tools/cache"
)

// GqlSourceLister helps list GqlSources.
type GqlSourceLister interface {
	// List lists all GqlSources in the indexer.
	List(selector labels.Selector) (ret []*v1alpha1.GqlSource, err error)
	// GqlSources returns an object that can list and get GqlSources.
	GqlSources(namespace string) GqlSourceNamespaceLister
	GqlSourceListerExpansion
}

// gqlSourceLister implements the GqlSourceLister interface.
type gqlSourceLister struct {
	indexer cache.Indexer
}

// NewGqlSourceLister returns a new GqlSourceLister.
func NewGqlSourceLister(indexer cache.Indexer) GqlSourceLister {
	return &gqlSourceLister{indexer: indexer}
}

// List lists all GqlSources in the indexer.
func (s *gqlSourceLister) List(selector labels.Selector) (ret []*v1alpha1.GqlSource, err error) {
	err = cache.ListAll(s.indexer, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GqlSource))
	})
	return ret, err
}

// GqlSources returns an object that can list and get GqlSources.
func (s *gqlSourceLister) GqlSources(namespace string) GqlSourceNamespaceLister {
	return gqlSourceNamespaceLister{indexer: s.indexer, namespace: namespace}
}

// GqlSourceNamespaceLister helps list and get GqlSources.
type GqlSourceNamespaceLister interface {
	// List lists all GqlSources in the indexer for a given namespace.
	List(selector labels.Selector) (ret []*v1alpha1.GqlSource, err error)
	// Get retrieves the GqlSource from the indexer for a given namespace and name.
	Get(name string) (*v1alpha1.GqlSource, error)
	GqlSourceNamespaceListerExpansion
}

// gqlSourceNamespaceLister implements the GqlSourceNamespaceLister
// interface.
type gqlSourceNamespaceLister struct {
	indexer   cache.Indexer
	namespace string
}

// List lists all GqlSources in the indexer for a given namespace.
func (s gqlSourceNamespaceLister) List(selector labels.Selector) (ret []*v1alpha1.GqlSource, err error) {
	err = cache.ListAllByNamespace(s.indexer, s.namespace, selector, func(m interface{}) {
		ret = append(ret, m.(*v1alpha1.GqlSource))
	})
	return ret, err
}

// Get retrieves the GqlSource from the indexer for a given namespace and name.
func (s gqlSourceNamespaceLister) Get(name string) (*v1alpha1.GqlSource, error) {
	obj, exists, err := s.indexer.GetByKey(s.namespace + "/" + name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, errors.NewNotFound(v1alpha1.Resource("gqlsource"), name)
	}
	return obj.(*v1alpha1.GqlSource), nil
}
