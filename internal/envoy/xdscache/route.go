package xdscache

import (
	route "github.com/envoyproxy/go-control-plane/envoy/config/route/v3"
	"sync"
)

type RouteCache struct {
	mu     sync.Mutex
	values map[string]*route.RouteConfiguration
}
