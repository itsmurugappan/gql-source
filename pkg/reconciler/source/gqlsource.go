package source

import (
	"context"
	"fmt"

	"go.uber.org/zap"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/equality"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	appsv1listers "k8s.io/client-go/listers/apps/v1"

	"knative.dev/eventing/pkg/reconciler/source"
	"knative.dev/pkg/apis"
	"knative.dev/pkg/logging"
	pkgreconciler "knative.dev/pkg/reconciler"
	"knative.dev/pkg/resolver"

	"github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
	"github.com/itsmurugappan/gql-source/pkg/client/clientset/versioned"
	reconcilergqlsource "github.com/itsmurugappan/gql-source/pkg/client/injection/reconciler/sources/v1alpha1/gqlsource"
	listers "github.com/itsmurugappan/gql-source/pkg/client/listers/sources/v1alpha1"
	"github.com/itsmurugappan/gql-source/pkg/reconciler/source/resources"
)

const (
	raImageEnvVar              = "GQL_RA_IMAGE"
	gqlSourceDeploymentCreated = "GqlSourceDeploymentCreated"
	gqlSourceDeploymentUpdated = "GqlSourceDeploymentUpdated"
	gqlSourceDeploymentFailed  = "GqlSourceDeploymentUpdated"
	component                  = "gqlsource"
)

// newDeploymentCreated makes a new reconciler event with event type Normal, and
// reason gqlSourceDeploymentCreated.
func newDeploymentCreated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, gqlSourceDeploymentCreated, "GqlSource created deployment: \"%s/%s\"", namespace, name)
}

// deploymentUpdated makes a new reconciler event with event type Normal, and
// reason gqlSourceDeploymentUpdated.
func deploymentUpdated(namespace, name string) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeNormal, gqlSourceDeploymentUpdated, "GqlSource updated deployment: \"%s/%s\"", namespace, name)
}

// newDeploymentFailed makes a new reconciler event with event type Warning, and
// reason gqlSourceDeploymentFailed.
func newDeploymentFailed(namespace, name string, err error) pkgreconciler.Event {
	return pkgreconciler.NewEvent(corev1.EventTypeWarning, gqlSourceDeploymentFailed, "GqlSource failed to create deployment: \"%s/%s\", %w", namespace, name, err)
}

type Reconciler struct {
	// KubeClientSet allows us to talk to the k8s for core APIs
	KubeClientSet kubernetes.Interface

	receiveAdapterImage string

	gqlLister        listers.GqlSourceLister
	deploymentLister appsv1listers.DeploymentLister

	gqlClientSet   versioned.Interface
	loggingContext context.Context

	sinkResolver *resolver.URIResolver

	configs source.ConfigAccessor
}

// Check that our Reconciler implements Interface
var _ reconcilergqlsource.Interface = (*Reconciler)(nil)

func (r *Reconciler) ReconcileKind(ctx context.Context, src *v1alpha1.GqlSource) pkgreconciler.Event {
	src.Status.InitializeConditions()

	if src.Spec.Sink == nil {
		src.Status.MarkNoSink("SinkMissing", "")
		return fmt.Errorf("spec.sink missing")
	}

	dest := src.Spec.Sink.DeepCopy()
	if dest.Ref != nil {
		// To call URIFromDestination(), dest.Ref must have a Namespace. If there is
		// no Namespace defined in dest.Ref, we will use the Namespace of the source
		// as the Namespace of dest.Ref.
		if dest.Ref.Namespace == "" {
			dest.Ref.Namespace = src.GetNamespace()
		}
	}
	sinkURI, err := r.sinkResolver.URIFromDestinationV1(*dest, src)
	if err != nil {
		src.Status.MarkNoSink("NotFound", "")
		return fmt.Errorf("getting sink URI: %v", err)
	}
	src.Status.MarkSink(sinkURI)

	ra, err := r.createReceiveAdapter(ctx, src, sinkURI)
	if err != nil {
		var event *pkgreconciler.ReconcilerEvent
		isReconcilerEvent := pkgreconciler.EventAs(err, &event)
		if isReconcilerEvent && event.EventType != corev1.EventTypeNormal {
			logging.FromContext(ctx).Error("Unable to create the receive adapter. Reconciler error", zap.Error(err))
			return err
		} else if !isReconcilerEvent {
			logging.FromContext(ctx).Error("Unable to create the receive adapter. Generic error", zap.Error(err))
			return err
		}
	}
	src.Status.MarkDeployed(ra)

	return nil
}

func (r *Reconciler) createReceiveAdapter(ctx context.Context, src *v1alpha1.GqlSource, sinkURI *apis.URL) (*appsv1.Deployment, error) {
	raArgs := resources.ReceiveAdapterArgs{
		Image:          r.receiveAdapterImage,
		Source:         src,
		Labels:         resources.GetLabels(src.Name),
		SinkURI:        sinkURI.String(),
		AdditionalEnvs: r.configs.ToEnvVars(),
	}
	expected := resources.MakeReceiveAdapter(&raArgs)

	ra, err := r.KubeClientSet.AppsV1().Deployments(src.Namespace).Get(expected.Name, metav1.GetOptions{})
	if err != nil && apierrors.IsNotFound(err) {
		ra, err = r.KubeClientSet.AppsV1().Deployments(src.Namespace).Create(expected)
		if err != nil {
			return nil, newDeploymentFailed(ra.Namespace, ra.Name, err)
		}
		return ra, newDeploymentCreated(ra.Namespace, ra.Name)
	} else if err != nil {
		logging.FromContext(ctx).Error("Unable to get an existing receive adapter", zap.Error(err))
		return nil, err
	} else if !metav1.IsControlledBy(ra, src) {
		return nil, fmt.Errorf("deployment %q is not owned by GqlSource %q", ra.Name, src.Name)
	} else if podSpecChanged(ra.Spec.Template.Spec, expected.Spec.Template.Spec) {
		ra.Spec.Template.Spec = expected.Spec.Template.Spec
		if ra, err = r.KubeClientSet.AppsV1().Deployments(src.Namespace).Update(ra); err != nil {
			return ra, err
		}
		return ra, deploymentUpdated(ra.Namespace, ra.Name)
	} else {
		logging.FromContext(ctx).Debug("Reusing existing receive adapter", zap.Any("receiveAdapter", ra))
	}
	return ra, nil
}

func podSpecChanged(oldPodSpec corev1.PodSpec, newPodSpec corev1.PodSpec) bool {
	if !equality.Semantic.DeepDerivative(newPodSpec, oldPodSpec) {
		return true
	}
	if len(oldPodSpec.Containers) != len(newPodSpec.Containers) {
		return true
	}
	for i := range newPodSpec.Containers {
		if !equality.Semantic.DeepEqual(newPodSpec.Containers[i].Env, oldPodSpec.Containers[i].Env) {
			return true
		}
	}
	return false
}
