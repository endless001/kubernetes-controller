package manager

import (
	"context"
	"fmt"
	"kubernetes-controller/internal/util/kubernetes/object/status"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
	"sigs.k8s.io/controller-runtime/pkg/manager"
)

func Run(ctx context.Context, c *Config) error {

	setupLog := ctrl.Log.WithName("setup")

	mgr, err := manager.New(config.GetConfigOrDie(), manager.Options{})
	if err != nil {
		return fmt.Errorf("unable to start controller manager: %w", err)
	}

	var kubernetesStatusQueue *status.Queue
	if c.UpdateStatus {
		setupLog.Info("Starting Status Updater")
		kubernetesStatusQueue = status.NewQueue()
	} else {
		setupLog.Info("status updates disabled, skipping status updater")
	}

	controllers, err := setupControllers(mgr, kubernetesStatusQueue, c)
	if err != nil {
		return fmt.Errorf("unable to setup controller as expected %w", err)
	}
	for _, c := range controllers {
		if err := c.MaybeSetupWithManager(mgr); err != nil {
			return fmt.Errorf("unable to create controller %q: %w", c.Name(), err)
		}
	}
	setupLog.Info("Starting manager")
	return mgr.Start(ctx)
}
