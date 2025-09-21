package tracker

import (
	"log"
	"sync"
	"time"

	"git-time-tracker.com/internal"
	"git-time-tracker.com/internal/commands"
)

type TrackData struct {
	mutex        sync.Mutex
	ProjectName  string
	ActiveBranch string
	ActiveTime   time.Duration
	LastTime     time.Time
}

var DATA = map[string]*TrackData{}

// This function is used to track the time spent on a project
// working by calculating active time since last changes in hashes or branch
// and saving it to the file after idle time or branch change
//
// @param projectName string - the name of the project
// @param branch string - the branch of the project
// @return error - the error if any
func Track(projectName string, branch string) error {
	data, ok := DATA[projectName]
	if !ok {
		data = &TrackData{
			mutex:        sync.Mutex{},
			ProjectName:  projectName,
			ActiveBranch: branch,
			ActiveTime:   time.Duration(0),
			LastTime:     time.Now(),
		}

		DATA[projectName] = data
		log.Printf("new project `%s` registered\n", projectName)
	}

	data.mutex.Lock()
	defer data.mutex.Unlock()

	// update after branch change
	if branch != data.ActiveBranch {
		log.Printf("branch changed on project `%s` from `%s` to `%s`\n", data.ProjectName, data.ActiveBranch, branch)
		saveActiveTime(data)

		// After calculate active time, reset all values
		data.ActiveBranch = branch
		data.ActiveTime = time.Duration(0)
		data.LastTime = time.Now()
		return nil
	}

	if internal.DEBUG {
		log.Printf(
			"add %s and sum %s to active time on project `%s`\n",
			time.Since(data.LastTime).Round(time.Second).String(),
			(data.ActiveTime + time.Since(data.LastTime)).Round(time.Second).String(),
			data.ProjectName,
		)
	}
	data.ActiveTime = data.ActiveTime + time.Since(data.LastTime)
	data.LastTime = time.Now()

	return nil
}

// This function is used to force save the active time to the file
// @param projectName string - the name of the project
// @return void
func ForceSave(projectName string) {
	data, ok := DATA[projectName]
	if !ok {
		return
	}

	data.mutex.Lock()
	defer data.mutex.Unlock()

	if data.ActiveTime <= 0 {
		if internal.DEBUG {
			log.Printf("no active time to save to project `%s`\n", data.ProjectName)
		}

		return
	}

	saveActiveTime(data)
	data.ActiveTime = time.Duration(0)
	data.LastTime = time.Now()
	return
}

// This function is used to save the active time to the file
// @param data *TrackData - the data of the project
// @return void
func saveActiveTime(data *TrackData) {
	if data.LastTime.Before(time.Now().Add(-internal.CHECK_INTERVAL)) {
		data.ActiveTime += internal.CHECK_INTERVAL
		if internal.DEBUG {
			log.Printf("add %s to active time on project `%s` on branch `%s`\n", internal.CHECK_INTERVAL.Round(time.Second).String(), data.ProjectName, data.ActiveBranch)
		}

	} else {
		data.ActiveTime += time.Since(data.LastTime)
		if internal.DEBUG {
			log.Printf("add %s to active time on project `%s` on branch `%s`\n", time.Since(data.LastTime).Round(time.Second).String(), data.ProjectName, data.ActiveBranch)
		}
	}

	log.Printf("saving active time to project `%s` on branch `%s`: %s\n", data.ProjectName, data.ActiveBranch, data.ActiveTime.Round(time.Second).String())

	copy := data.ActiveTime
	err := commands.AddTimeToProjectAndBranch(data.ProjectName, data.ActiveBranch, &copy)
	if err != nil {
		log.Printf("error adding time to project and branch: %v", err)
	}
}
