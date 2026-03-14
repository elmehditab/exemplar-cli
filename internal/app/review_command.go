package app

import (
	"errors"
	"os"
)

type ReviewRequest struct {
	RepoPath string
}

func RunReview(req ReviewRequest) (string, error) {
	if req.RepoPath == "" {
		return "", errors.New("repo path is required")
	}
	info, err := os.Stat(req.RepoPath)

	if err != nil {
		return "", errors.New("repo path does not exist: " + req.RepoPath)
	}
	if !info.IsDir() {
		return "", errors.New("repo path must be a directory")
	}

	return "review command invoked for repo: " + req.RepoPath, nil
}
