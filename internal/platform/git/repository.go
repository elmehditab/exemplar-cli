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
