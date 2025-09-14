package commands

import (
	"crypto/md5"
	"encoding/hex"
	"os/exec"
	"strings"
)

func GetGitStatusHash(repoPath string) (string, error) {
	cmd := exec.Command("git", "status", "--short")
	cmd.Dir = repoPath
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	hash := md5.Sum(output)
	return hex.EncodeToString(hash[:]), nil
}

func GetGitBranch(repoPath string) (string, error) {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	output, err := cmd.Output()

	if err != nil {
		return "", err
	}

	branch := strings.ReplaceAll(string(output), "\n", "")

	return branch, nil
}
