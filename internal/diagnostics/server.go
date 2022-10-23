package diagnostics

import (
	"github.com/go-logr/logr"
	"sync"
)

type Server struct {
	Logger           logr.Logger
	ProfilingEnabled bool
	ConfigLock       *sync.RWMutex
}
