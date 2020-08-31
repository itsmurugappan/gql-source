package main

import (
	"knative.dev/eventing/pkg/adapter/v2"

	gqladapter "github.com/itsmurugappan/gql-source/pkg/adapter"
)

func main() {
	adapter.Main("gqlsource", gqladapter.NewEnvConfig, gqladapter.NewAdapter)
}
