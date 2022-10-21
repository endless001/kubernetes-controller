package root

import (
	"context"
	"errors"
	"kubernetes-controller/internal/manager"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var (
	mutex           sync.Mutex
	shutdownSignals = []os.Signal{os.Interrupt, syscall.SIGTERM}
)

func SetupSignalHandler(cfg *manager.Config) (context.Context, error) {
	if ok := mutex.TryLock(); !ok {
		return nil, errors.New("signal handler can only be setup once")
	}
	ctx, cancel := context.WithCancel(context.Background())

	c := make(chan os.Signal, 2)
	signal.Notify(c, shutdownSignals...)
	go func() {
		<-c
		select {
		case <-time.After(cfg.TermDelay):
			cancel()
		case <-c:
			os.Exit(1) // second signal. Exit directly.
		}

		<-c
		os.Exit(1) // second signal. Exit directly.
	}()

	return ctx, nil
}
