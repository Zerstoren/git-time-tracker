package commands

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"time"

	"git-time-tracker.com/internal"
)

const format = "%s | %s : %s"
const searchFormat = "%s | %s :"

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

func writeTimeToFile(projectName string, branch string, duration *time.Duration) error {
	fileInfo := internal.FILE_PATH
	file, err := os.OpenFile(fileInfo, os.O_RDWR|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	lines := []string{}

	scanner := bufio.NewScanner(file)

	isExist := false
	for scanner.Scan() {
		line := scanner.Text()

		if strings.Contains(line, fmt.Sprintf(searchFormat, projectName, branch)) {
			lines = append(lines, fmt.Sprintf(format, projectName, branch, duration.Round(time.Second).String()))
			isExist = true
		} else {
			lines = append(lines, line)
		}
	}

	if !isExist {
		lines = append(lines, fmt.Sprintf(format, projectName, branch, duration.Round(time.Second).String()))
	}

	err = os.WriteFile(fileInfo, []byte(strings.Join(lines, "\n")), 0644)
	if err != nil {
		return err
	}

	return nil
}
