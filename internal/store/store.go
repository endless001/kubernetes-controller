package store

import (
	"encoding/json"
	"fmt"
	corev1 "k8s.io/api/core/v1"
	extensions "k8s.io/api/extensions/v1beta1"
	netv1 "k8s.io/api/networking/v1"
	netv1beta1 "k8s.io/api/networking/v1beta1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/tools/cache"
	"reflect"
	"sort"
	"strings"
	"sync"
)

const (
	IngressClassController = "inendless.com/ingress-controller"
)

type ErrNotFound struct {
	message string
}

func (e ErrNotFound) Error() string {
	if e.message == "" {
		return "not found"
	}
	return e.message
}

type Storer interface {
	GetSecret(namespace, name string) (*corev1.Secret, error)
	GetService(namespace, name string) (*corev1.Service, error)
	GetEndpointsForService(namespace, name string) (*corev1.Endpoints, error)
	GetIngressClassV1(name string) (*netv1.IngressClass, error)
	ListIngressesV1beta1() []*netv1beta1.Ingress
	ListIngressesV1() []*netv1.Ingress
	ListIngressClassesV1() []*netv1.IngressClass
}

type Store struct {
	stores       CacheStores
	ingressClass string
}

type CacheStores struct {
	IngressV1beta1 cache.Store
	IngressV1      cache.Store
	IngressClassV1 cache.Store
	Service        cache.Store
	Secret         cache.Store
	Endpoint       cache.Store

	l *sync.RWMutex
}

func NewCacheStores() CacheStores {
	return CacheStores{
		IngressV1beta1: cache.NewStore(keyFunc),
		IngressV1:      cache.NewStore(keyFunc),
		IngressClassV1: cache.NewStore(clusterResourceKeyFunc),
		Service:        cache.NewStore(keyFunc),
		Secret:         cache.NewStore(keyFunc),
		Endpoint:       cache.NewStore(keyFunc),

		l: &sync.RWMutex{},
	}
}

func (c CacheStores) Get(obj runtime.Object) (item interface{}, exists bool, err error) {
	c.l.RLock()
	defer c.l.RUnlock()
	switch obj := obj.(type) {
	case *extensions.Ingress:
		return c.IngressV1beta1.Get(obj)
	case *netv1beta1.Ingress:
		return c.IngressV1beta1.Get(obj)
	case *netv1.Ingress:
		return c.IngressV1.Get(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Get(obj)
	case *corev1.Service:
		return c.Service.Get(obj)
	case *corev1.Secret:
		return c.Secret.Get(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Get(obj)
	default:
		return nil, false, fmt.Errorf("%T is not a supported cache object type", obj)
	}
}

func (c CacheStores) Add(obj runtime.Object) error {
	c.l.Lock()
	defer c.l.Unlock()

	switch obj := obj.(type) {
	case *extensions.Ingress:
		return c.IngressV1beta1.Add(obj)
	case *netv1beta1.Ingress:
		return c.IngressV1beta1.Add(obj)
	case *netv1.Ingress:
		return c.IngressV1.Add(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Add(obj)
	case *corev1.Service:
		return c.Service.Add(obj)
	case *corev1.Secret:
		return c.Secret.Add(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Add(obj)
	default:
		return fmt.Errorf("cannot add unsupported kind %q to the store", obj.GetObjectKind().GroupVersionKind())
	}
}

func (c CacheStores) Delete(obj runtime.Object) error {
	c.l.Lock()
	defer c.l.Unlock()
	switch obj := obj.(type) {
	case *extensions.Ingress:
		return c.IngressV1beta1.Delete(obj)
	case *netv1beta1.Ingress:
		return c.IngressV1beta1.Delete(obj)
	case *netv1.Ingress:
		return c.IngressV1.Delete(obj)
	case *netv1.IngressClass:
		return c.IngressClassV1.Delete(obj)
	case *corev1.Service:
		return c.Service.Delete(obj)
	case *corev1.Secret:
		return c.Secret.Delete(obj)
	case *corev1.Endpoints:
		return c.Endpoint.Delete(obj)
	default:
		return fmt.Errorf("cannot delete unsupported kind %q from the store", obj.GetObjectKind().GroupVersionKind())

	}
}

func New(cs CacheStores, ingressClass string) Storer {
	return Store{
		stores:       cs,
		ingressClass: ingressClass,
	}
}

func (s Store) GetSecret(namespace, name string) (*corev1.Secret, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	secret, exists, err := s.stores.Secret.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Secret %v not found", key)}
	}
	return secret.(*corev1.Secret), nil
}
func (s Store) GetService(namespace, name string) (*corev1.Service, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	service, exists, err := s.stores.Service.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Service %v not found", key)}
	}
	return service.(*corev1.Service), nil
}

func (s Store) ListIngressesV1() []*netv1.Ingress {
	var ingresses []*netv1.Ingress
	for _, item := range s.stores.IngressV1.List() {
		ing, ok := item.(*netv1.Ingress)
		if !ok {
			continue
		}
		ingresses = append(ingresses, ing)
	}
	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})

	return ingresses
}

func (s Store) ListIngressClassesV1() []*netv1.IngressClass {
	var classes []*netv1.IngressClass
	for _, item := range s.stores.IngressClassV1.List() {
		class, ok := item.(*netv1.IngressClass)
		if !ok {
			continue
		}
		if class.Spec.Controller != IngressClassController {
			continue
		}
		classes = append(classes, class)
	}
	sort.SliceStable(classes, func(i, j int) bool {
		return strings.Compare(classes[i].Name, classes[j].Name) < 0
	})
	return classes
}

func (s Store) GetEndpointsForService(namespace, name string) (*corev1.Endpoints, error) {
	key := fmt.Sprintf("%v/%v", namespace, name)
	eps, exists, err := s.stores.Endpoint.GetByKey(key)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("Endpoints for service %v not found", key)}
	}
	return eps.(*corev1.Endpoints), nil
}

func (s Store) ListIngressesV1beta1() []*netv1beta1.Ingress {
	var ingresses []*netv1beta1.Ingress
	for _, item := range s.stores.IngressV1beta1.List() {
		ing := s.networkingIngressV1Beta1(item)
		ingresses = append(ingresses, ing)
	}
	sort.SliceStable(ingresses, func(i, j int) bool {
		return strings.Compare(fmt.Sprintf("%s/%s", ingresses[i].Namespace, ingresses[i].Name),
			fmt.Sprintf("%s/%s", ingresses[j].Namespace, ingresses[j].Name)) < 0
	})
	return ingresses
}
func (s Store) GetIngressClassV1(name string) (*netv1.IngressClass, error) {
	p, exists, err := s.stores.IngressClassV1.GetByKey(name)
	if err != nil {
		return nil, err
	}
	if !exists {
		return nil, ErrNotFound{fmt.Sprintf("IngressClass %v not found", name)}
	}
	return p.(*netv1.IngressClass), nil
}

func (s Store) networkingIngressV1Beta1(obj interface{}) *netv1beta1.Ingress {
	switch obj := obj.(type) {
	case *netv1beta1.Ingress:
		return obj

	case *extensions.Ingress:
		out, err := toNetworkingIngressV1Beta1(obj)
		if err != nil {
			return nil
		}
		return out

	default:
		return nil
	}
}

func toNetworkingIngressV1Beta1(obj *extensions.Ingress) (*netv1beta1.Ingress, error) {
	js, err := json.Marshal(obj)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize object of type %v: %w", reflect.TypeOf(obj), err)
	}
	var out netv1beta1.Ingress
	if err := json.Unmarshal(js, &out); err != nil {
		return nil, fmt.Errorf("failed to deserialize json: %w", err)
	}
	out.APIVersion = netv1beta1.SchemeGroupVersion.String()
	return &out, nil
}
func keyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	name := v.FieldByName("Name")
	namespace := v.FieldByName("Namespace")
	return namespace.String() + "/" + name.String(), nil
}

func clusterResourceKeyFunc(obj interface{}) (string, error) {
	v := reflect.Indirect(reflect.ValueOf(obj))
	return v.FieldByName("Name").String(), nil
}
