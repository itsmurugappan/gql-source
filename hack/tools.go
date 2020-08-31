// +build tools

package tools

// This package imports things required by this repository, to force `go mod` to see them as dependencies
import (
	_ "k8s.io/code-generator"
	_ "knative.dev/test-infra/scripts"

	// codegen: hack/generate-knative.sh
	_ "knative.dev/pkg/hack"

	_ "k8s.io/code-generator/cmd/client-gen"
	_ "k8s.io/code-generator/cmd/deepcopy-gen"
	_ "k8s.io/code-generator/cmd/defaulter-gen"
	_ "k8s.io/code-generator/cmd/informer-gen"
	_ "k8s.io/code-generator/cmd/lister-gen"
	_ "k8s.io/kube-openapi/cmd/openapi-gen"
	_ "knative.dev/pkg/codegen/cmd/injection-gen"

	// Licenseclassifier
	_ "github.com/google/licenseclassifier"
)
