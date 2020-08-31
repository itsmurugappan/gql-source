module github.com/itsmurugappan/gql-source

go 1.14

require (
	github.com/cloudevents/sdk-go/v2 v2.2.0
	github.com/google/licenseclassifier v0.0.0-20200708223521-3d09a0ea2f39
	github.com/google/uuid v1.1.1
	github.com/itsmurugappan/graphql v0.0.0-00010101000000-000000000000
	github.com/matryer/is v1.4.0 // indirect
	go.uber.org/zap v1.15.0
	k8s.io/api v0.18.7-rc.0
	k8s.io/apimachinery v0.18.7-rc.0
	k8s.io/client-go v11.0.1-0.20190805182717-6502b5e7b1b5+incompatible
	k8s.io/code-generator v0.18.6
	k8s.io/kube-openapi v0.0.0-20200410145947-bcb3869e6f29
	knative.dev/eventing v0.17.1
	knative.dev/pkg v0.0.0-20200819202314-b5411f2221aa
	knative.dev/test-infra v0.0.0-20200819210814-f578ab25945b
)

replace (
	github.com/itsmurugappan/graphql => github.com/itsmurugappan/graphql v0.2.3-0.20200214202050-b9117788f6e3
	github.com/prometheus/client_golang => github.com/prometheus/client_golang v0.9.2
	k8s.io/api => k8s.io/api v0.17.6
	k8s.io/apiextensions-apiserver => k8s.io/apiextensions-apiserver v0.17.6
	k8s.io/apimachinery => k8s.io/apimachinery v0.17.6
	k8s.io/client-go => k8s.io/client-go v0.17.6
	k8s.io/code-generator => k8s.io/code-generator v0.17.6
)
