package tracker

import (
	"fmt"
	"log"
	"sync"
	"time"

	"git-time-tracker.com/internal/commands"
)

type TrackData struct {
	ProjectName  string
	Paths        []string
	LastHashes   []string
	ActiveBranch string
	ActiveTime   time.Duration
	LastTime     time.Time
}

var DATA = map[string]*TrackData{}

func Track(wg *sync.WaitGroup, projectName string, paths []string) error {
	wg.Add(1)
	defer wg.Done()

	// Get branch name only from first path
	branch, err := commands.GetGitBranch(paths[0])
	if err != nil {
		return err
	}

	hashes, err := getHash(paths)
	if err != nil {
		return err
	}

	data, ok := DATA[projectName]
	if !ok {
		data = &TrackData{
			ProjectName:  projectName,
			Paths:        paths,
			LastHashes:   hashes,
			ActiveBranch: branch,
			ActiveTime:   time.Duration(0),
			LastTime:     time.Now(),
		}

		DATA[projectName] = data
	}

	// update after idle time
	// if data.LastTime.Before(time.Now().Add(-internal.MAX_IDLE_TIME)) {
	// 	fmt.Println("idle time")
	// 	saveActiveTime(data)
	// 	data.ActiveTime = time.Duration(0)
	// 	data.LastTime = time.Now()
	// }

	// update after branch change
	if branch != data.ActiveBranch {
		saveActiveTime(data)
		data.LastHashes = hashes
		data.ActiveBranch = branch
		data.ActiveTime = time.Duration(0)
		data.LastTime = time.Now()
	}

	// update after hashes changed asda dsffsdfs
	if isHashesChanged(data.LastHashes, hashes) {
		data.LastHashes = hashes
		data.ActiveTime += time.Since(data.LastTime)
		data.LastTime = time.Now()
		fmt.Println("hashes changed", data.ActiveTime)
	}

	return nil
}

func getHash(paths []string) ([]string, error) {
	hashes := []string{}
	for _, path := range paths {
		hash, err := commands.GetGitStatusHash(path)
		if err != nil {
			return nil, err
		}
		hashes = append(hashes, hash)
	}

	return hashes, nil
}

func isHashesChanged(lastHashes []string, newHashes []string) bool {
	for i, _ := range lastHashes {
		fmt.Println("lastHashes[i]", lastHashes[i])
		fmt.Println("newHashes[i]", newHashes[i])
		if newHashes[i] != lastHashes[i] {
			return true
		}
	}
	return false
}

func saveActiveTime(data *TrackData) {
	copy := data.ActiveTime
	err := commands.AddTimeToProjectAndBranch(data.ProjectName, data.ActiveBranch, &copy)
	if err != nil {
		log.Printf("error adding time to project and branch: %v", err)
	} else {
		log.Printf("time added to project %s and branch %s: %s", data.ProjectName, data.ActiveBranch, data.ActiveTime.String())
	}
}
