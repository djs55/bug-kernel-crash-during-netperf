package main

import (
	"log"
	"time"

	"github.com/pkg/errors"
)

// RestartDesktop reboots the VM which should fix problems where a FUSE container deadlocks during
// `docker kill` (probably because the server has gone away unexpectedly)
func RestartDesktop() error {
	guiClient := NewClient()
	initialState, err := guiClient.GetLastEvent()
	if err != nil {
		return err
	}
	if err := guiClient.Restart(); err != nil {
		return err
	}
	for {
		time.Sleep(time.Second)
		event, err := NewClient().GetLastEvent()
		if err != nil {
			log.Printf("waiting for a connection to Desktop: %v", err)
			continue
		}
		if event.State == FailedToStart {
			return errors.New("failed to start Docker")
		}
		if event.State != Running {
			log.Println("waiting for Desktop to get into a Running state")
			continue
		}
		if event.Timestamp <= initialState.Timestamp {
			log.Println("waiting for Desktop to restart and report a new timestamp")
			continue
		}
		return nil
	}
}
