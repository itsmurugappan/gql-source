package v1alpha1

import (
	appsv1 "k8s.io/api/apps/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"knative.dev/eventing/pkg/apis/duck"
	"knative.dev/pkg/apis"
)

const (
	// GqlSourceConditionSinkProvided has status True when the GqlSource has been configured with a sink target.
	GqlSourceConditionSinkProvided apis.ConditionType = "SinkProvided"

	// GqlSourceConditionDeployed has status True when the GqlSource has had it's receive adapter deployment created.
	GqlSourceConditionDeployed apis.ConditionType = "Deployed"

	// GqlSourceConditionReady has status True when the GqlSource is ready to send events.
	GqlSourceConditionReady = apis.ConditionReady
)

var condSet = apis.NewLivingConditionSet(
	GqlSourceConditionSinkProvided,
	GqlSourceConditionDeployed,
)

// GetGroupVersionKind implements kmeta.OwnerRefable
func (g *GqlSource) GetGroupVersionKind() schema.GroupVersionKind {
	return SchemeGroupVersion.WithKind("GqlSource")
}

// GetConditionSet retrieves the condition set for this resource. Implements the KRShaped interface.
func (g *GqlSource) GetConditionSet() apis.ConditionSet {
	return condSet
}

func (s *GqlSourceStatus) GetCondition(t apis.ConditionType) *apis.Condition {
	return condSet.Manage(s).GetCondition(t)
}

// IsReady returns true if the resource is ready overall.
func (s *GqlSourceStatus) IsReady() bool {
	return condSet.Manage(s).IsHappy()
}

// InitializeConditions sets the initial values to the conditions.
func (gs *GqlSourceStatus) InitializeConditions() {
	condSet.Manage(gs).InitializeConditions()
}

// MarkSink sets the condition that the source has a sink configured.
func (gs *GqlSourceStatus) MarkSink(uri *apis.URL) {
	gs.SinkURI = uri
	if !uri.IsEmpty() {
		condSet.Manage(gs).MarkTrue(GqlSourceConditionSinkProvided)
	} else {
		condSet.Manage(gs).MarkUnknown(GqlSourceConditionSinkProvided, "SinkEmpty", "Sink has resolved to empty.%s", "")
	}
}

// MarkNoSink sets the condition that the source does not have a sink configured.
func (gs *GqlSourceStatus) MarkNoSink(reason, messageFormat string, messageA ...interface{}) {
	condSet.Manage(gs).MarkFalse(GqlSourceConditionSinkProvided, reason, messageFormat, messageA...)
}

func DeploymentIsAvailable(d *appsv1.DeploymentStatus, def bool) bool {
	// Check if the Deployment is available.
	for _, cond := range d.Conditions {
		if cond.Type == appsv1.DeploymentAvailable {
			return cond.Status == "True"
		}
	}
	return def
}

// MarkDeployed sets the condition that the source has been deployed.
func (s *GqlSourceStatus) MarkDeployed(d *appsv1.Deployment) {
	if duck.DeploymentIsAvailable(&d.Status, false) {
		condSet.Manage(s).MarkTrue(GqlSourceConditionDeployed)
	} else {
		// I don't know how to propagate the status well, so just give the name of the Deployment
		// for now.
		condSet.Manage(s).MarkFalse(GqlSourceConditionDeployed, "DeploymentUnavailable", "The Deployment '%s' is unavailable.", d.Name)
	}
}

// MarkDeploying sets the condition that the source is deploying.
func (s *GqlSourceStatus) MarkDeploying(reason, messageFormat string, messageA ...interface{}) {
	condSet.Manage(s).MarkUnknown(GqlSourceConditionDeployed, reason, messageFormat, messageA...)
}

// MarkNotDeployed sets the condition that the source has not been deployed.
func (s *GqlSourceStatus) MarkNotDeployed(reason, messageFormat string, messageA ...interface{}) {
	condSet.Manage(s).MarkFalse(GqlSourceConditionDeployed, reason, messageFormat, messageA...)
}
