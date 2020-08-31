package resources

const (
	// controllerAgentName is the string used by this controller to identify
	// itself when creating events.
	controllerAgentName = "gql-source-controller"
)

func GetLabels(name string) map[string]string {
	return map[string]string{
		"eventing.knative.dev/source":     controllerAgentName,
		"eventing.knative.dev/SourceName": name,
	}
}
