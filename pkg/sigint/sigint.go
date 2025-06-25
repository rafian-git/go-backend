package sigint

import (
	"os"
	"os/signal"
	"syscall"
)

// Wait for SIGINT, SIGTERM or SIGQUIT and return.
func Wait() {
	var sig = make(chan os.Signal, 2)
	signal.Notify(sig, os.Interrupt, syscall.SIGQUIT, syscall.SIGTERM)
	<-sig
}
