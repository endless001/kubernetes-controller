package manager

import (
	"context"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	serverv3 "github.com/envoyproxy/go-control-plane/pkg/server/v3"
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	"kubernetes-controller/internal/envoy/xds"
	ctrl "sigs.k8s.io/controller-runtime"
)

func setupLoggers(c *Config) (logrus.FieldLogger, logr.Logger, error) {
	panic("")
}

func setupControllerOptions(logger logr.Logger, c *Config, scheme *runtime.Scheme) (ctrl.Options, error) {

	var leaderElection bool

	controllerOpts := ctrl.Options{
		Scheme:                 scheme,
		MetricsBindAddress:     c.MetricsAddr,
		Port:                   9443,
		HealthProbeBindAddress: c.ProbeAddr,
		LeaderElection:         leaderElection,
		LeaderElectionID:       c.LeaderElectionID,
		SyncPeriod:             &c.SyncPeriod,
	}
	return controllerOpts, nil
}

func setupXdsServer(ctx context.Context) error {
	cache := cache.NewSnapshotCache(false, cache.IDHash{}, nil)
	srv := serverv3.NewServer(ctx, cache, nil)
	_ = xds.NewServer(srv)
	return nil
}
