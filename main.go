package main

import (
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"git-time-tracker.com/internal"
	"git-time-tracker.com/internal/tracker"
)

func main() {
	internal.ReadConfig()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool, 1)

	go func() {
	loop:
		for {
			wg := sync.WaitGroup{}
			for project, paths := range internal.REPOSITORIES {
				go tracker.Track(&wg, project, paths)
			}

			select {
			case <-sigs:
				break loop
			case <-time.After(internal.CHECK_INTERVAL):
				wg.Wait()
			}
		}

		fmt.Println("Program terminated. Save data to file.")

		for project, paths := range internal.REPOSITORIES {
			tracker.ForceSave(project, paths)
		}

		done <- true
	}()

	fmt.Println("Program running. Press Ctrl+C or send SIGTERM to stop.")

	<-done
	fmt.Println("Program terminated.")
}

// asdas 111112111111
