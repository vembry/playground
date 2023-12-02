package common

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

// WatchForExitSignal is to awaits incoming interrupt signal
// sent to the service
func WatchForExitSignal() os.Signal {
	log.Printf("awaiting sigterm...")
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
