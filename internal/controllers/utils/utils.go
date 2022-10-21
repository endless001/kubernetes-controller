package utils

import (
	netv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/meta"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"kubernetes-controller/internal/annotations"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/event"
	"sigs.k8s.io/controller-runtime/pkg/predicate"
)

const defaultIngressClassAnnotation = "ingressclass.kubernetes.io/is-default-class"

func IsDefaultIngressClass(obj client.Object) bool {
	if ingressClass, ok := obj.(*netv1.IngressClass); ok {
		return ingressClass.ObjectMeta.Annotations[defaultIngressClassAnnotation] == "true"
	}
	return false
}
func GeneratePredicateFuncsForIngressClassFilter(name string) predicate.Funcs {
	preds := predicate.NewPredicateFuncs(func(obj client.Object) bool {
		return MatchesIngressClass(obj, name, true)
	})
	preds.UpdateFunc = func(e event.UpdateEvent) bool {
		return MatchesIngressClass(e.ObjectOld, name, true) || MatchesIngressClass(e.ObjectNew, name, true)
	}
	return preds
}

func MatchesIngressClass(obj client.Object, controllerIngressClass string, isDefault bool) bool {
	objectIngressClass := obj.GetAnnotations()[annotations.IngressClassKey]
	if isDefault && IsIngressClassEmpty(obj) {
		return true
	}
	if ing, isV1Ingress := obj.(*netv1.Ingress); isV1Ingress {
		if ing.Spec.IngressClassName != nil && *ing.Spec.IngressClassName == controllerIngressClass {
			return true
		}
	}

	switch controllerIngressClass {
	case objectIngressClass:
		return true
	}
	return false
}

func IsIngressClassEmpty(obj client.Object) bool {
	switch obj := obj.(type) {
	case *netv1.Ingress:
		if _, ok := obj.GetAnnotations()[annotations.IngressClassKey]; !ok {
			return obj.Spec.IngressClassName == nil
		}
		return false
	default:
		if _, ok := obj.GetAnnotations()[annotations.IngressClassKey]; ok {
			return false
		}
		return true
	}
}
func CRDExists(restMapper meta.RESTMapper, gvr schema.GroupVersionResource) bool {
	_, err := restMapper.KindFor(gvr)
	return !meta.IsNoMatchError(err)
}
