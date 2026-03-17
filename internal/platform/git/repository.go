package git

import (
	"errors"
	"os/exec"

	"strings"
)

func ResolveGitRepository(repoPath string) (string, error) {

	cmd := exec.Command("git", "-C", repoPath, "rev-parse", "--show-toplevel")
	out, err := cmd.Output()

	if err != nil {
		return "", errors.New("repo path is not a git repository: " + repoPath)
	}

	gitRepo := strings.TrimSpace(string(out))

	return gitRepo, nil
}

func GetCurrentBranch(repoPath string) (string, error) {

	cmd := exec.Command("git", "-C", repoPath, "branch", "--show-current")
	out, err := cmd.Output()

	if err != nil {
		return "", errors.New("failed to get current branch for repository: " + repoPath)
	}

	branch := strings.TrimSpace(string(out))

	return branch, nil
}
