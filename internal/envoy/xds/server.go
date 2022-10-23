package xds

import (
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"google.golang.org/grpc"
)

type Server interface {
	clusterservice.ClusterDiscoveryServiceServer
	endpointservice.EndpointDiscoveryServiceServer
	listenerservice.ListenerDiscoveryServiceServer
	routeservice.RouteDiscoveryServiceServer
	discoverygrpc.AggregatedDiscoveryServiceServer
	secretservice.SecretDiscoveryServiceServer
	runtimeservice.RuntimeDiscoveryServiceServer
}

func RegisterServer(srv Server, g *grpc.Server) {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(g, srv)
	secretservice.RegisterSecretDiscoveryServiceServer(g, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(g, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(g, srv)
	listenerservice.RegisterListenerDiscoveryServiceServer(g, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(g, srv)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(g, srv)
}
