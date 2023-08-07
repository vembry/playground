package cmd

import (
	"os"
	"os/signal"
	"syscall"
)

// watchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func watchForExitSignal() os.Signal {
	ch := make(chan os.Signal, 4)
	signal.Notify(
		ch,
		syscall.SIGINT,
		syscall.SIGQUIT,
		syscall.SIGTERM,
		syscall.SIGTSTP,
	)

	return <-ch
}
