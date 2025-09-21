package internal

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Repository struct {
	Path    string   `json:"path"`
	Exclude []string `json:"exclude"`
}

type Config struct {
	Mode             *string               `json:"mode"`
	CheckInterval    string                `json:"check_interval"`
	DebounceInterval string                `json:"debounce_interval"`
	MaxIdleTime      int                   `json:"max_idle_time"`
	WriteToFile      bool                  `json:"write_to_file"`
	FilePath         string                `json:"file_path"`
	LogFilePath      string                `json:"log_file_path"`
	Repositories     map[string]Repository `json:"repositories"`
}

var DEBUG = false
var CHECK_INTERVAL = 20 * time.Minute
var DEBOUNCE_INTERVAL = 2 * time.Second
var WRITE_TO_FILE = true
var FILE_PATH = "time-tracker.txt"
var LOG_FILE_PATH = "time-logs.log"
var REPOSITORIES = map[string]Repository{}

// This function is used to read the config file
// @return void
func ReadConfig() {
	config := &Config{}

	configFile, err := os.ReadFile("config.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(configFile, &config)
	if err != nil {
		panic(err)
	}

	if config.Mode != nil {
		DEBUG = *config.Mode == "debug"
	}

	for project, paths := range config.Repositories {
		if strings.Contains(project, ":") {
			// stop program with fatal error
			panic("Project name can`t contain ':'")
		}

		REPOSITORIES[project] = paths
	}

	if config.LogFilePath != "" {
		write, err := os.OpenFile(config.LogFilePath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			panic(err)
		}

		log.SetOutput(write)
	}

	checkInterval, err := time.ParseDuration(config.CheckInterval)
	if err != nil {
		panic(err)
	}

	debounceInterval, err := time.ParseDuration(config.DebounceInterval)
	if err != nil {
		panic(err)
	}

	CHECK_INTERVAL = checkInterval
	WRITE_TO_FILE = config.WriteToFile
	FILE_PATH = config.FilePath
	REPOSITORIES = config.Repositories
	DEBOUNCE_INTERVAL = debounceInterval
}
