package commands

import (
	"crypto/md5"
	"encoding/hex"
	"os"
	"os/exec"
	"strings"
)

// This function is used to get the hash of the git status
// @param repoPath string - the path of the repository
// @return string - the hash of the git status
// @return error - the error if execution failed
func GetGitStatusHash(repoPath string) (string, error) {
	cmd := exec.Command("git", "diff", "HEAD")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	hash := md5.Sum(output)
	return hex.EncodeToString(hash[:]), nil
}

// This function is used to get the branch of the git repository
// @param repoPath string - the path of the repository
// @return string - the branch of the git repository
// @return error - the error if execution failed
func GetGitBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	branch := strings.ReplaceAll(string(output), "\n", "")

	return branch, nil
}
