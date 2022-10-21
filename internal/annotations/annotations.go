package annotations

type ClassMatching int

const (
	IgnoreClassMatch       ClassMatching = iota
	ExactOrEmptyClassMatch ClassMatching = iota
	ExactClassMatch        ClassMatching = iota
)

const (
	AnnotationPrefix    = "inendless.com"
	IngressClassKey     = "kubernetes.io/ingress.class"
	DefaultIngressClass = "sail"
)
