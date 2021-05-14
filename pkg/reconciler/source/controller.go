package source

import (
	"context"
	"os"

	"k8s.io/client-go/tools/cache"

	"knative.dev/eventing/pkg/reconciler/source"

	kubeclient "knative.dev/pkg/client/injection/kube/client"
	deploymentinformer "knative.dev/pkg/client/injection/kube/informers/apps/v1/deployment"
	"knative.dev/pkg/configmap"
	"knative.dev/pkg/controller"
	"knative.dev/pkg/logging"
	"knative.dev/pkg/resolver"

	"github.com/itsmurugappan/gql-source/pkg/apis/sources/v1alpha1"
	gqlclient "github.com/itsmurugappan/gql-source/pkg/client/injection/client"
	gqlinformer "github.com/itsmurugappan/gql-source/pkg/client/injection/informers/sources/v1alpha1/gqlsource"
	"github.com/itsmurugappan/gql-source/pkg/client/injection/reconciler/sources/v1alpha1/gqlsource"
)

func NewController(
	ctx context.Context,
	cmw configmap.Watcher,
) *controller.Impl {

	raImage, defined := os.LookupEnv(raImageEnvVar)
	if !defined {
		logging.FromContext(ctx).Errorf("required environment variable '%s' not defined", raImageEnvVar)
		return nil
	}

	gqlinformer := gqlinformer.Get(ctx)
	deploymentInformer := deploymentinformer.Get(ctx)

	c := &Reconciler{
		KubeClientSet:       kubeclient.Get(ctx),
		gqlClientSet:        gqlclient.Get(ctx),
		gqlLister:           gqlinformer.Lister(),
		deploymentLister:    deploymentInformer.Lister(),
		receiveAdapterImage: raImage,
		loggingContext:      ctx,
		configs:             source.WatchConfigurations(ctx, component, cmw),
	}

	impl := gqlsource.NewImpl(ctx, c)
	c.sinkResolver = resolver.NewURIResolver(ctx, impl.EnqueueKey)

	logging.FromContext(ctx).Info("Setting up gql subscription event handlers")

	gqlinformer.Informer().AddEventHandler(controller.HandleAll(impl.Enqueue))

	deploymentInformer.Informer().AddEventHandler(cache.FilteringResourceEventHandler{
		FilterFunc: controller.FilterControllerGK(v1alpha1.Kind("GqlSource")),
		Handler:    controller.HandleAll(impl.EnqueueControllerOf),
	})

	return impl
}
