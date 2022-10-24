package configuration

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrlutils "kubernetes-controller/internal/controllers/utils"
	"kubernetes-controller/internal/store"
	"kubernetes-controller/internal/util"
	"kubernetes-controller/internal/util/kubernetes/object/status"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
	"time"
)

type CoreV1ServiceReconciler struct {
	client.Client
	cache            *store.CacheStores
	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration
}

func (r *CoreV1ServiceReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("CoreV1Service", mgr, controller.Options{
		Reconciler: r,

		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		&source.Kind{Type: &corev1.Service{}},
		&handler.EnqueueRequestForObject{},
	)
}

func (r *CoreV1ServiceReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Service", req.NamespacedName)
	obj := new(corev1.Service)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.cache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Service", "namespace", req.Namespace, "name", req.Name)

		_, objectExistsInCache, err := r.cache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.cache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
	if err := r.cache.Add(obj.DeepCopyObject()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

type CoreV1EndpointsReconciler struct {
	client.Client
	cache            *store.CacheStores
	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration
}

func (r *CoreV1EndpointsReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("CoreV1Endpoints", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		&source.Kind{Type: &corev1.Endpoints{}},
		&handler.EnqueueRequestForObject{},
	)
}

func (r *CoreV1EndpointsReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {

	log := r.Log.WithValues("CoreV1Endpoints", req.NamespacedName)
	// get the relevant object
	obj := new(corev1.Endpoints)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.cache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Endpoints", "namespace", req.Namespace, "name", req.Name)
		_, objectExistsInCache, err := r.cache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.cache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
	if err := r.cache.Add(obj.DeepCopyObject()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

type CoreV1SecretReconciler struct {
	client.Client
	cache            *store.CacheStores
	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration
}

func (r *CoreV1SecretReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("CoreV1Secret", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		&source.Kind{Type: &corev1.Secret{}},
		&handler.EnqueueRequestForObject{},
	)
}
func (r *CoreV1SecretReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("CoreV1Secret", req.NamespacedName)
	obj := new(corev1.Secret)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.cache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Secret", "namespace", req.Namespace, "name", req.Name)
		_, objectExistsInCache, err := r.cache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.cache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
	if err := r.cache.Add(obj.DeepCopyObject()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

type NetV1IngressReconciler struct {
	client.Client

	Log                        logr.Logger
	Scheme                     *runtime.Scheme
	CacheSyncTimeout           time.Duration
	Cache                      *store.CacheStores
	StatusQueue                *status.Queue
	IngressClassName           string
	DisableIngressClassLookups bool
}

func (r *NetV1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("NetV1Ingress", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "networking.k8s.io",
				Version: "v1",
				Kind:    "Ingress",
			})},
			&handler.EnqueueRequestForObject{},
		); err != nil {
			return err
		}
	}
	if !r.DisableIngressClassLookups {
		err = c.Watch(
			&source.Kind{Type: &netv1.IngressClass{}},
			handler.EnqueueRequestsFromMapFunc(r.listClassless),
			predicate.NewPredicateFuncs(ctrlutils.IsDefaultIngressClass),
		)
		if err != nil {
			return err
		}
	}
	preds := ctrlutils.GeneratePredicateFuncsForIngressClassFilter(r.IngressClassName)
	return c.Watch(
		&source.Kind{Type: &netv1.Ingress{}},
		&handler.EnqueueRequestForObject{},
		preds,
	)
}

func (r *NetV1IngressReconciler) listClassless(obj client.Object) []reconcile.Request {
	resourceList := &netv1.IngressList{}
	if err := r.Client.List(context.Background(), resourceList); err != nil {
		r.Log.Error(err, "failed to list classless ingresses")
		return nil
	}
	var recs []reconcile.Request
	for i, resource := range resourceList.Items {
		if ctrlutils.IsIngressClassEmpty(&resourceList.Items[i]) {
			recs = append(recs, reconcile.Request{
				NamespacedName: types.NamespacedName{
					Namespace: resource.Namespace,
					Name:      resource.Name,
				},
			})
		}
	}
	return recs
}
func (r *NetV1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1Ingress", req.NamespacedName)

	// get the relevant object
	obj := new(netv1.Ingress)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.Cache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "Ingress", "namespace", req.Namespace, "name", req.Name)
		_, objectExistsInCache, err := r.Cache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.Cache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}

	class := new(netv1.IngressClass)
	if !r.DisableIngressClassLookups {
		if err := r.Get(ctx, types.NamespacedName{Name: r.IngressClassName}, class); err != nil {
			// we log this without taking action to support legacy configurations that only set ingressClassName or
			// used the class annotation and did not create a corresponding IngressClass. We only need this to determine
			// if the IngressClass is default or to configure default settings, and can assume no/no additional defaults
			// if none exists.
			log.V(util.DebugLevel).Info("could not retrieve IngressClass", "ingressclass", r.IngressClassName)
		}
	}
	// if the object is not configured with our ingress.class, then we need to ensure it's removed from the cache
	if !ctrlutils.MatchesIngressClass(obj, r.IngressClassName, ctrlutils.IsDefaultIngressClass(class)) {
		log.V(util.DebugLevel).Info("object missing ingress class, ensuring it's removed from configuration",
			"namespace", req.Namespace, "name", req.Name, "class", r.IngressClassName)
		return ctrl.Result{}, r.Cache.Delete(obj)
	} else {
		log.V(util.DebugLevel).Info("object has matching ingress class", "namespace", req.Namespace, "name", req.Name,
			"class", r.IngressClassName)
	}
	if err := r.Cache.Add(obj.DeepCopyObject()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}

type NetV1IngressClassReconciler struct {
	client.Client
	Cache            *store.CacheStores
	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration
}

func (r *NetV1IngressClassReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("NetV1IngressClass", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	return c.Watch(
		&source.Kind{Type: &netv1.IngressClass{}},
		&handler.EnqueueRequestForObject{},
	)
}

func (r *NetV1IngressClassReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	log := r.Log.WithValues("NetV1IngressClass", req.NamespacedName)

	// get the relevant object
	obj := new(netv1.IngressClass)
	if err := r.Get(ctx, req.NamespacedName, obj); err != nil {
		if errors.IsNotFound(err) {
			obj.Namespace = req.Namespace
			obj.Name = req.Name
			return ctrl.Result{}, r.Cache.Delete(obj)
		}
		return ctrl.Result{}, err
	}
	log.V(util.DebugLevel).Info("reconciling resource", "namespace", req.Namespace, "name", req.Name)

	// clean the object up if it's being deleted
	if !obj.DeletionTimestamp.IsZero() && time.Now().After(obj.DeletionTimestamp.Time) {
		log.V(util.DebugLevel).Info("resource is being deleted, its configuration will be removed", "type", "IngressClass", "namespace", req.Namespace, "name", req.Name)
		_, objectExistsInCache, err := r.Cache.Get(obj)
		if err != nil {
			return ctrl.Result{}, err
		}
		if objectExistsInCache {
			if err := r.Cache.Delete(obj); err != nil {
				return ctrl.Result{}, err
			}
			return ctrl.Result{Requeue: true}, nil // wait until the object is no longer present in the cache
		}
		return ctrl.Result{}, nil
	}
	if err := r.Cache.Add(obj.DeepCopyObject()); err != nil {
		return ctrl.Result{}, err
	}
	return ctrl.Result{}, nil
}
