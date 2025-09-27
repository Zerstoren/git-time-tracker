# Git Time Tracker

A powerful Go-based time tracking tool that automatically monitors your development activity across multiple Git repositories. It tracks time spent on different branches and projects by monitoring file system changes and Git branch switches.

## Features

- **Automatic Time Tracking**: Monitors file changes in your Git repositories to track active development time
- **Branch-Aware**: Tracks time separately for each Git branch
- **Multi-Repository Support**: Monitor multiple projects simultaneously
- **File System Monitoring**: Uses `fsnotify` to detect file changes in real-time
- **Debounced Events**: Prevents excessive tracking from rapid file changes
- **Configurable**: Highly customizable through JSON configuration
- **Persistent Storage**: Saves tracked time to a text file for analysis
- **Debug Mode**: Detailed logging for troubleshooting

## How It Works

The tracker monitors file system events in your configured Git repositories. When files are modified, created, or deleted, it:

1. Detects the current Git branch
2. Tracks time spent on that branch
3. Automatically saves time when switching branches or after idle periods
4. Stores the data in a persistent format

## Installation

### Prerequisites

- Go 1.25.0 or later
- Git repositories to track

### Build from Source

```bash
git clone https://github.com/your-username/git-time-tracker.git
cd git-time-tracker
go mod tidy
go build -o git-time-tracker
```

## Configuration

Create a `config.json` file at executable file:

```json
{
  "mode": "debug",
  "check_interval": "10m",
  "debounce_interval": "2s",
  "write_to_file": true,
  
  "file_path": "../time-tracker.txt",
  "log_file_path": "../time-logs.log",

  "repositories": {
    "my-project": {
      "path": "/path/to/your/repository",
      "exclude": [
        "/node_modules",
        "/dist",
        "/.git",
        "/build"
      ]
    },
    "another-project": {
      "path": "/path/to/another/repository",
      "exclude": [
        "/vendor",
        "/.git"
      ]
    }
  }
}
```

### Configuration Options

- **mode**: `"debug"` or `"production"` - Controls logging verbosity
- **check_interval**: How often to save tracked time (e.g., `"10m"`, `"30m"`)
- **debounce_interval**: Delay before processing file changes (e.g., `"2s"`, `"5s"`)
- **write_to_file**: Whether to save tracked time to a file
- **file_path**: Path to the output file for tracked time data
- **log_file_path**: Path to the log file (optional)
- **repositories**: Object mapping project names to repository configurations
  - **path**: Absolute path to the Git repository
  - **exclude**: Array of paths to exclude from monitoring

## Usage

1. Configure your repositories in `config.json`
2. Run the tracker:

```bash
./git-time-tracker
```

The program will:
- Start monitoring all configured repositories
- Display startup messages
- Run continuously until interrupted with `Ctrl+C`

### Output Format

Tracked time is saved in the following format:
```
project-name | branch-name : 2h30m15s
my-project | main : 1h45m30s
my-project | feature-branch : 45m20s
```

## Dependencies

- `github.com/fsnotify/fsnotify` - File system event monitoring
- `golang.org/x/sys` - System-specific functionality

## Troubleshooting

### Common Issues

1. **Permission Errors**: Ensure the application has read/write access to repository paths and output files
2. **Git Branch Detection Fails**: Verify that the repository path is correct and contains a valid Git repository
3. **File Monitoring Issues**: Check that excluded paths are correctly configured

### Debug Information

Enable debug mode to see detailed information about:
- File system events
- Branch changes
- Time calculations
- File operations

## License

This project is open source. Please check the license file for details.
