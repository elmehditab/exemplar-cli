package git

import (
	"bytes"
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

	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain", "--untracked-files=all")
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

	untrackedDiff, err := getUntrackedDiff(repoPath)
	if err != nil {
		return "", err
	}

	return diff + untrackedDiff, nil
}

func getUntrackedDiff(repoPath string) (string, error) {
	files, err := getUntrackedFiles(repoPath)
	if err != nil {
		return "", err
	}

	var diff bytes.Buffer

	for _, file := range files {
		out, err := runNoIndexDiff(repoPath, file)
		if err != nil {
			return "", err
		}
		diff.Write(out)
	}

	return diff.String(), nil
}

func getUntrackedFiles(repoPath string) ([]string, error) {
	cmd := exec.Command("git", "-C", repoPath, "status", "--porcelain", "--untracked-files=all")
	out, err := cmd.Output()

	if err != nil {
		return nil, errors.New("failed to get untracked files for repository: " + repoPath)
	}

	if len(out) == 0 {
		return []string{}, nil
	}

	lines := strings.Split(string(out), "\n")
	var untrackedFiles []string

	for _, line := range lines {
		if strings.HasPrefix(line, "?? ") && len(line) > 3 {
			untrackedFiles = append(untrackedFiles, strings.TrimSpace(line[3:]))
		}
	}

	return untrackedFiles, nil
}

func runNoIndexDiff(repoPath string, file string) ([]byte, error) {
	cmd := exec.Command("git", "-C", repoPath, "diff", "--no-ext-diff", "--no-index", "--", "/dev/null", file)
	out, err := cmd.Output()

	if err == nil {
		return out, nil
	}

	if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
		return out, nil
	}

	return nil, errors.New("failed to get diff for untracked file: " + file)
}
