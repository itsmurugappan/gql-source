package main

import (
	// The set of controllers this controller process runs.
	"github.com/itsmurugappan/gql-source/pkg/reconciler/source"

	// This defines the shared main for injected controllers.
	"knative.dev/pkg/injection/sharedmain"
)

func main() {
	sharedmain.Main("controller",
		source.NewController,
	)
}
