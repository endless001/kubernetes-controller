package xds

import (
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"google.golang.org/grpc"
)

func registerServer(g *grpc.Server, srv serverv3.Server) {
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(g, srv)
	secretservice.RegisterSecretDiscoveryServiceServer(g, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(g, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(g, srv)
	listenerservice.RegisterListenerDiscoveryServiceServer(g, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(g, srv)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(g, srv)
}

func NewServer(server serverv3.Server, opts ...grpc.ServerOption) *grpc.Server {
	g := grpc.NewServer(opts...)
	registerServer(g, server)
	return g
}
