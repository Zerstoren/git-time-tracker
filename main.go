package main

import (
	"fmt"
	"sync"
	"time"

	"git-time-tracker.com/internal"
	"git-time-tracker.com/internal/tracker"
)

// test comment asdas sadas
// qweqw
func main() {
	internal.ReadConfig()

	for {
		wg := sync.WaitGroup{}
		for project, paths := range internal.REPOSITORIES {
			go tracker.Track(&wg, project, paths)
		}

		wg.Wait()

		time.Sleep(internal.CHECK_INTERVAL)
		fmt.Println("check interval")
	}
}
