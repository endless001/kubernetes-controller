package configuration

import (
	"context"
	"github.com/go-logr/logr"
	corev1 "k8s.io/api/core/v1"
	extv1beta1 "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	ctrlutils "kubernetes-controller/internal/controllers/utils"
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
	panic("")
}

type CoreV1EndpointsReconciler struct {
	client.Client

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
	panic("")
}

type CoreV1SecretReconciler struct {
	client.Client

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
	panic("")
}

type NetV1IngressReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration

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
	panic("1")
}

type NetV1IngressClassReconciler struct {
	client.Client

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
	panic("")
}

type NetV1Beta1IngressReconciler struct {
	client.Client

	Log              logr.Logger
	Scheme           *runtime.Scheme
	CacheSyncTimeout time.Duration

	StatusQueue                *status.Queue
	IngressClassName           string
	DisableIngressClassLookups bool
}

func (r *NetV1Beta1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("NetV1Beta1Ingress", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "networking.k8s.io",
				Version: "v1beta1",
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
		&source.Kind{Type: &netv1beta1.Ingress{}},
		&handler.EnqueueRequestForObject{},
		preds,
	)
}
func (r *NetV1Beta1IngressReconciler) listClassless(obj client.Object) []reconcile.Request {
	resourceList := &netv1beta1.IngressList{}
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
func (r *NetV1Beta1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	panic("")
}

type ExtV1Beta1IngressReconciler struct {
	client.Client

	Log                        logr.Logger
	Scheme                     *runtime.Scheme
	CacheSyncTimeout           time.Duration
	StatusQueue                *status.Queue
	IngressClassName           string
	DisableIngressClassLookups bool
}

func (r *ExtV1Beta1IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	c, err := controller.New("ExtV1Beta1Ingress", mgr, controller.Options{
		Reconciler: r,
		LogConstructor: func(_ *reconcile.Request) logr.Logger {
			return r.Log
		},
		CacheSyncTimeout: r.CacheSyncTimeout,
	})
	if err != nil {
		return err
	}
	// if configured, start the status updater controller
	if r.StatusQueue != nil {
		if err := c.Watch(
			&source.Channel{Source: r.StatusQueue.Subscribe(schema.GroupVersionKind{
				Group:   "extensions",
				Version: "v1beta1",
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
		&source.Kind{Type: &extv1beta1.Ingress{}},
		&handler.EnqueueRequestForObject{},
		preds,
	)
}
func (r *ExtV1Beta1IngressReconciler) listClassless(obj client.Object) []reconcile.Request {
	resourceList := &extv1beta1.IngressList{}
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
func (r *ExtV1Beta1IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	panic("")
}
