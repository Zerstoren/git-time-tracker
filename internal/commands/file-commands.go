package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"git-time-tracker.com/internal"
)

const format = "%s | %s : %s"
const searchFormat = "%s | %s :"

var mutex = sync.Mutex{}

// This function is used to add the time to the project and branch
// @param projectName string - the name of the project
// @param branch string - the branch of the project
// @param durationAdd *time.Duration - the time to add
// @return error - the error if any
func AddTimeToProjectAndBranch(projectName string, branch string, durationAdd *time.Duration) error {
	durationFromFile, err := getTimeFrom(projectName, branch)
	if err != nil {
		return err
	}

	duration := *durationFromFile + *durationAdd

	err = writeTimeToFile(projectName, branch, &duration)
	if err != nil {
		return err
	}

	return nil
}

// This function is used to get the time from the file
// probably can be slow, but it's not a big deal at now
// better solution is to use database like sqlite
//
// @param projectName string - the name of the project
// @param branch string - the branch of the project
// @return *time.Duration - the time from the file
// @return error - the error if any
func getTimeFrom(projectName string, branch string) (*time.Duration, error) {
	fileInfo := internal.FILE_PATH
	file, err := os.Open(fileInfo)
	if err != nil {
		timeDuration := time.Duration(0)
		return &timeDuration, nil
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.Contains(line, fmt.Sprintf(searchFormat, projectName, branch)) {
			continue
		}

		timePeriod := strings.Split(line, ":")[1]
		duration, err := time.ParseDuration(strings.TrimSpace(timePeriod))

		if err != nil {
			return nil, err
		}

		return &duration, nil
	}

	timeDuration := time.Duration(0)
	return &timeDuration, nil
}

// This function is used to write the time to the file
// @param projectName string - the name of the project
// @param branch string - the branch of the project
// @param duration *time.Duration - the time to write
// @return error - the error if any
func writeTimeToFile(projectName string, branch string, duration *time.Duration) error {

	// Prevent race condition
	mutex.Lock()
	defer mutex.Unlock()

	fileInfo := internal.FILE_PATH
	file, err := os.OpenFile(fileInfo, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}

	scanner := bufio.NewScanner(file)

	// Read file line by line
	isExist := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, fmt.Sprintf(searchFormat, projectName, branch)) {
			// When line matches pattern, update it
			lines = append(lines, fmt.Sprintf(format, projectName, branch, duration.Round(time.Second).String()))
			isExist = true
		} else {
			// When line doesn't match pattern, add it to the list
			lines = append(lines, line)
		}
	}

	// When line doesn't exist, add it to the list
	if !isExist {
		lines = append(lines, fmt.Sprintf(format, projectName, branch, duration.Round(time.Second).String()))
	}

	err = os.WriteFile(fileInfo, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}
