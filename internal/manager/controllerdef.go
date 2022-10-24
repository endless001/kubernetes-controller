package manager

import (
	"fmt"
	"kubernetes-controller/internal/controllers/configuration"
	"kubernetes-controller/internal/store"
	"kubernetes-controller/internal/util/kubernetes/object/status"
	"reflect"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

type Controller interface {
	SetupWithManager(ctrl.Manager) error
}

type ControllerDef struct {
	Enabled    bool
	Controller Controller
}

func (c *ControllerDef) Name() string {
	return reflect.TypeOf(c.Controller).String()
}

func (c *ControllerDef) MaybeSetupWithManager(mgr ctrl.Manager) error {
	if !c.Enabled {
		return nil
	}
	return c.Controller.SetupWithManager(mgr)
}

func setupControllers(
	mgr manager.Manager,
	kubernetesStatusQueue *status.Queue,
	cache *store.CacheStores,
	c *Config) ([]ControllerDef, error) {

	restMapper := mgr.GetClient().RESTMapper()
	ingressConditions, err := NewIngressControllersConditions(c, restMapper)
	if err != nil {
		return nil, fmt.Errorf("ingress version picker failed: %w", err)
	}

	controllers := []ControllerDef{
		{
			Enabled: ingressConditions.IngressClassNetV1Enabled(),
			Controller: &configuration.NetV1IngressClassReconciler{
				Client:           mgr.GetClient(),
				Cache:            cache,
				Log:              ctrl.Log.WithName("controllers").WithName("IngressClass").WithName("netv1"),
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: ingressConditions.IngressNetV1Enabled(),
			Controller: &configuration.NetV1IngressReconciler{
				Client:                     mgr.GetClient(),
				Cache:                      cache,
				Log:                        ctrl.Log.WithName("controllers").WithName("Ingress").WithName("netv1"),
				Scheme:                     mgr.GetScheme(),
				IngressClassName:           c.IngressClassName,
				DisableIngressClassLookups: !c.IngressClassNetV1Enabled,
				StatusQueue:                kubernetesStatusQueue,
				CacheSyncTimeout:           c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1ServiceReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Service"),
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: c.ServiceEnabled,
			Controller: &configuration.CoreV1EndpointsReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Endpoints"),
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
		{
			Enabled: true,
			Controller: &configuration.CoreV1SecretReconciler{
				Client:           mgr.GetClient(),
				Log:              ctrl.Log.WithName("controllers").WithName("Secrets"),
				Scheme:           mgr.GetScheme(),
				CacheSyncTimeout: c.CacheSyncTimeout,
			},
		},
	}
	return controllers, nil
}
