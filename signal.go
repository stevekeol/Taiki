package main

import (
	"os"
	"os/signal"
)

var (
	shutdownRequestChannel = make(chan struct{})
	interruptSignals       = []os.Signal{os.Interrupt}
)

func interruptListeners() <-chan struct{} {
	c := make(chan struct{})
	go func() {
		interruptChannel := make(chan os.Signal, 1)
		signal.Notify(interruptChannel, interruptSignals...)

		// Listen for initial shutdown signal and close the returned
		// channel to notify the caller.
		select {
		case sig := <-interruptChannel:
			log.Debug("Received signal then Shutting down...", "sig", sig)

		case <-shutdownRequestChannel:
			log.Debug("Shutdown requested then Shutting down...")
		}
		close(c)

		// Listen for repeated signals and display a message so the user
		// knows the shutdown is in progress and the process is not
		// hung.
		for {
			select {
			case <-interruptChannel:
				log.Debug("Received signal. Already shutting down...")

			case <-shutdownRequestChannel:
				log.Debug("Shutdown requested. Already shutting down...")
			}
		}
	}()
	return c
}

// interruptRequested returns true when the channel returned by
// interruptListener was closed.  This simplifies early shutdown slightly since
// the caller can just use an if statement instead of a select.
func interruptRequested(interrupted <-chan struct{}) bool {
	select {
	case <-interrupted:
		return true
	default:
	}

	return false
}
