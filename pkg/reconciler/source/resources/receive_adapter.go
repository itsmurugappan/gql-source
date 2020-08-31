package resources

import (
	"fmt"
	"strings"

	"github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
	v1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"knative.dev/eventing/pkg/utils"
	"knative.dev/pkg/kmeta"
)

type ReceiveAdapterArgs struct {
	Image          string
	Source         *v1alpha1.GqlSource
	Labels         map[string]string
	SinkURI        string
	AdditionalEnvs []corev1.EnvVar
}

func MakeReceiveAdapter(args *ReceiveAdapterArgs) *v1.Deployment {
	replicas := int32(1)

	env := append([]corev1.EnvVar{{
		Name:  "GQL_SERVER",
		Value: args.Source.Spec.GqlServer,
	}, {
		Name:  "SUBSCRIPTION_QUERIES",
		Value: strings.Join(args.Source.Spec.SubscriptionQueries, ","),
	}, {
		Name:  "K_SINK",
		Value: args.SinkURI,
	}, {
		Name:  "NAME",
		Value: args.Source.Name,
	}, {
		Name:  "NAMESPACE",
		Value: args.Source.Namespace,
	}}, args.AdditionalEnvs...)

	return &v1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      utils.GenerateFixedName(args.Source, fmt.Sprintf("gqlsource-%s", args.Source.Name)),
			Namespace: args.Source.Namespace,
			Labels:    args.Labels,
			OwnerReferences: []metav1.OwnerReference{
				*kmeta.NewControllerRef(args.Source),
			},
		},
		Spec: v1.DeploymentSpec{
			Selector: &metav1.LabelSelector{
				MatchLabels: args.Labels,
			},
			Replicas: &replicas,
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Annotations: map[string]string{
						"sidecar.istio.io/inject": "false",
					},
					Labels: args.Labels,
				},
				Spec: corev1.PodSpec{
					// ServiceAccountName: args.Source.Spec.ServiceAccountName,
					Containers: []corev1.Container{
						{
							Name:  "receive-adapter",
							Image: args.Image,
							Env:   env,
						},
					},
				},
			},
		},
	}
}
