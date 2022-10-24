package xdscache

import (
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"kubernetes-controller/internal/envoy/resources"
)

type Cache struct {
	Listeners map[string]resources.Listener
	Routes    map[string]resources.Route
	Clusters  map[string]resources.Cluster
	Endpoints map[string]resources.Endpoint
}

func (cache *Cache) ClusterContents() []types.Resource {
	var r []types.Resource

	for _, c := range cache.Clusters {
		r = append(r, resources.MakeCluster(c.Name))
	}

	return r
}

func (cache *Cache) RouteContents() []types.Resource {

	var routesArray []resources.Route
	for _, r := range cache.Routes {
		routesArray = append(routesArray, r)
	}

	return []types.Resource{resources.MakeRoute(routesArray)}
}

func (cache *Cache) ListenerContents() []types.Resource {
	var r []types.Resource

	for _, l := range cache.Listeners {
		r = append(r, resources.MakeHTTPListener(l.Name, l.RouteNames[0], l.Address, l.Port))
	}

	return r
}

func (cache *Cache) EndpointsContents() []types.Resource {
	var r []types.Resource

	for _, c := range cache.Clusters {
		r = append(r, resources.MakeEndpoint(c.Name, c.Endpoints))
	}

	return r
}
func (cache *Cache) AddListener(name string, routeNames []string, address string, port uint32) {
	cache.Listeners[name] = resources.Listener{
		Name:       name,
		Address:    address,
		Port:       port,
		RouteNames: routeNames,
	}
}
func (cache *Cache) AddRoute(name, prefix string, clusters []string) {
	cache.Routes[name] = resources.Route{
		Name:    name,
		Prefix:  prefix,
		Cluster: clusters[0],
	}
}
func (cache *Cache) AddCluster(name string) {
	cache.Clusters[name] = resources.Cluster{
		Name: name,
	}
}
func (cache *Cache) AddEndpoint(clusterName, upstreamHost string, upstreamPort uint32) {
	cluster := cache.Clusters[clusterName]

	cluster.Endpoints = append(cluster.Endpoints, resources.Endpoint{
		UpstreamHost: upstreamHost,
		UpstreamPort: upstreamPort,
	})

	cache.Clusters[clusterName] = cluster
}
