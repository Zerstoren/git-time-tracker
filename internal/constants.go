package internal

import (
	"encoding/json"
	"log"
	"os"
	"strings"
	"time"
)

type Config struct {
	CheckInterval int                 `json:"check_interval"`
	MaxIdleTime   int                 `json:"max_idle_time"`
	WriteToFile   bool                `json:"write_to_file"`
	FilePath      string              `json:"file_path"`
	LogFilePath   string              `json:"log_file_path"`
	Repositories  map[string][]string `json:"repositories"`
}

const DEFAULT_SIZE = time.Minute

var CHECK_INTERVAL = 20 * DEFAULT_SIZE
var WRITE_TO_FILE = true
var FILE_PATH = "time-tracker.txt"
var LOG_FILE_PATH = "time-logs.log"
var REPOSITORIES = map[string][]string{}

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

	CHECK_INTERVAL = time.Duration(config.CheckInterval) * DEFAULT_SIZE
	WRITE_TO_FILE = config.WriteToFile
	FILE_PATH = config.FilePath
	REPOSITORIES = config.Repositories
}
