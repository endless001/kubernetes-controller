package root

import (
	"fmt"
	"kubernetes-controller/internal/manager"
)

func Run(c *manager.Config) error {
	ctx, err := SetupSignalHandler(c)
	if err != nil {
		return fmt.Errorf("failed to setup signal handler: %w", err)
	}
	return manager.Run(ctx, c)
}
