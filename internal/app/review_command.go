package app

import "errors"

type ReviewRequest struct {
	RepoPath string
}

func RunReview(req ReviewRequest) (string, error) {
	if req.RepoPath == "" {
		return "", errors.New("repo path is required")
	}
	return "review command invoked for repo: " + req.RepoPath, nil
}
