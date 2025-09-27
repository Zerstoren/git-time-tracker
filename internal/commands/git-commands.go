package commands

import (
	"os"
	"os/exec"
	"strings"
)

// This function is used to get the branch of the git repository
// @param repoPath string - the path of the repository
// @return string - the branch of the git repository
// @return error - the error if execution failed
func GetGitBranch(repoPath string) string {
	cmd := exec.Command("git", "branch", "--show-current")
	cmd.Dir = repoPath
	cmd.Stderr = os.Stderr
	output, err := cmd.Output()

	if err != nil {
		panic("error getting git branch")
	}

	branch := strings.ReplaceAll(string(output), "\n", "")

	return branch
}
