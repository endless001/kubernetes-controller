package envoy

import "google.golang.org/grpc"

func NewServer(opts ...grpc.ServerOption) *grpc.Server {
	g := grpc.NewServer(opts...)
	return g
}
