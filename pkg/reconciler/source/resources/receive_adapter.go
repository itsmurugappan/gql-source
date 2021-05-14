package resources

import (
	"fmt"
	"strconv"
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
		Name:  "GQL_OIDC_AUTH_ENABLE",
		Value: strconv.FormatBool(args.Source.Spec.OIDC.Enable),
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

	env = appendEnvFromSecretKeyRef(env, "GQL_AUTH_TOKEN_URL", args.Source.Spec.OIDC.TokenUrl.SecretKeyRef)
	env = appendEnvFromSecretKeyRef(env, "GQL_OIDC_CLIENT", args.Source.Spec.OIDC.OIDCClient.SecretKeyRef)
	env = appendEnvFromSecretKeyRef(env, "GQL_OIDC_USER", args.Source.Spec.OIDC.User.SecretKeyRef)
	env = appendEnvFromSecretKeyRef(env, "GQL_OIDC_PASSWORD", args.Source.Spec.OIDC.Password.SecretKeyRef)

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

// appendEnvFromSecretKeyRef returns env with an EnvVar appended
// setting key to the secret and key described by ref.
// If ref is nil, env is returned unchanged.
func appendEnvFromSecretKeyRef(env []corev1.EnvVar, key string, ref *corev1.SecretKeySelector) []corev1.EnvVar {
	if ref == nil {
		return env
	}

	env = append(env, corev1.EnvVar{
		Name: key,
		ValueFrom: &corev1.EnvVarSource{
			SecretKeyRef: ref,
		},
	})

	return env
}
