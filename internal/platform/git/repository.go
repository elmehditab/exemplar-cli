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

func GetChangedFiles(repoPath string) ([]string, error) {

	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain")
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.New("failed to get changed files for repository: " + repoPath)
	}

	if len(out) == 0 {
		return []string{}, nil
	}

	lines := strings.Split(string(out), "\n")
	var changedFiles []string

	for _, line := range lines {
		if len(line) > 3 {
			changedFiles = append(changedFiles, strings.TrimSpace(line[3:]))
		}
	}

	return changedFiles, nil
}

func GetDiff(repoPath string) (string, error) {

	cmd := exec.Command("git", "-C", repoPath, "diff", "--no-ext-diff", "HEAD")
	out, err := cmd.Output()

	if err != nil {
		return "", errors.New("failed to get diff for repository: " + repoPath)
	}

	diff := string(out)
	return diff, nil
}
