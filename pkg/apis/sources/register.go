package sources

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
)

const (
	GroupName = "sources.muru.dev"
)

var (
	GqlSourceResource = schema.GroupResource{
		Group:    GroupName,
		Resource: "gqlsources",
	}
)
