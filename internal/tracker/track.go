package tracker

import (
	"log"
	"sync"
	"time"

	"git-time-tracker.com/internal"
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

// This function is used to track the time spent on a project
// @param wg *sync.WaitGroup - the wait group to wait for the function to complete
// @param projectName string - the name of the project
// @param paths []string - the paths of the project
// @return error - the error if any
func Track(wg *sync.WaitGroup, projectName string, paths []string) error {
	wg.Add(1)
	defer wg.Done()

	// Get branch name only from first path
	branch, err := commands.GetGitBranch(paths[0])
	if err != nil {
		log.Printf("error getting branch: %v", err)
		return err
	}

	hashes, err := getHash(paths)
	if err != nil {
		log.Printf("error getting hash: %v", err)
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
	if data.LastTime.Before(time.Now().Add(-internal.MAX_IDLE_TIME)) {
		saveActiveTime(data)
		data.ActiveTime = time.Duration(0)
		data.LastTime = time.Now()
	}

	// update after branch change
	if branch != data.ActiveBranch {
		log.Printf("branch changed on project `%s` from `%s` to `%s`\n", data.ProjectName, data.ActiveBranch, branch)
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
		log.Printf("hashes changed on project `%s` on branch `%s`: %s\n", data.ProjectName, data.ActiveBranch, data.ActiveTime.Round(time.Second).String())
	}

	return nil
}

func ForceSave(projectName string, paths []string) {
	data, ok := DATA[projectName]
	if !ok {
		return
	}

	saveActiveTime(data)
	return
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
		if newHashes[i] != lastHashes[i] {
			return true
		}
	}
	return false
}

func saveActiveTime(data *TrackData) {
	if data.ActiveTime == 0 {
		log.Printf("no active time to save to project `%s` on branch `%s`\n", data.ProjectName, data.ActiveBranch)
		return
	}

	log.Printf("saving active time to project `%s` on branch `%s`: %s\n", data.ProjectName, data.ActiveBranch, data.ActiveTime.Round(time.Second).String())

	copy := data.ActiveTime
	err := commands.AddTimeToProjectAndBranch(data.ProjectName, data.ActiveBranch, &copy)
	if err != nil {
		log.Printf("error adding time to project and branch: %v", err)
	}
}
