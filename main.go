package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

	"git-time-tracker.com/internal"
	"git-time-tracker.com/internal/commands"
	"git-time-tracker.com/internal/tracker"
	"github.com/fsnotify/fsnotify"
)

func main() {
	internal.ReadConfig()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)

	done := make(chan bool, 1)
	wg := sync.WaitGroup{}
	wg.Add(1)

	for projectCycle, repositoryCycle := range internal.REPOSITORIES {
		go (func(project string, repository internal.Repository) {
			timer := time.NewTimer(internal.CHECK_INTERVAL)

			if internal.DEBUG {
				log.Printf("check interval on project `%s` every %s\n", project, internal.CHECK_INTERVAL.String())
			}

			go tracker.Notify(&wg, repository, func(event fsnotify.Event) {
				branch := commands.GetGitBranch(repository.Path)
				tracker.Track(project, branch)
				timer.Reset(internal.CHECK_INTERVAL)
			})

			for {
				select {
				case <-timer.C:
					if internal.DEBUG {
						log.Printf("check interval on project `%s`\n", project)
					}

					tracker.ForceSave(project)
					timer.Reset(internal.CHECK_INTERVAL)
				}
			}
		})(projectCycle, repositoryCycle)
	}

	go func() {
		<-sigs
		wg.Done()

		fmt.Println("Program terminated. Save data to file.")

		for project, _ := range internal.REPOSITORIES {
			tracker.ForceSave(project)
		}

		done <- true
	}()

	fmt.Println("Program running. Press Ctrl+C or send SIGTERM to stop.")

	<-done
	fmt.Println("Program terminated.")
}
