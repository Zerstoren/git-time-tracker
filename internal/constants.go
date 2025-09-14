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

var CHECK_INTERVAL = 60 * time.Second
var MAX_IDLE_TIME = 60 * time.Minute
var WRITE_TO_FILE = true
var FILE_PATH = "time-tracker.txt"
var LOG_FILE_PATH = "time-logs.log"
var REPOSITORIES = map[string][]string{}

func ReadConfig() *Config {
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
			panic("Project name cannot contain ':'")
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

	CHECK_INTERVAL = time.Duration(config.CheckInterval) * time.Second
	MAX_IDLE_TIME = time.Duration(config.MaxIdleTime) * time.Second
	WRITE_TO_FILE = config.WriteToFile
	FILE_PATH = config.FilePath
	REPOSITORIES = config.Repositories

	return config
}
