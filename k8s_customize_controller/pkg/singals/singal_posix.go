package singals

import (
	"os"
	"syscall"
)

var shutdownSingals = []os.Signal{os.Interrupt, syscall.SIGTERM}
