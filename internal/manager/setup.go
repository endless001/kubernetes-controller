package manager

import (
	"github.com/go-logr/logr"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
)

func setupLoggers(c *Config) (logrus.FieldLogger, logr.Logger, error) {
	panic("")
}

func setupControllerOptions(logger logr.Logger, c *Config, scheme *runtime.Scheme,
	dbMode string) (ctrl.Options, error) {

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
