package app

type ReviewRequest struct {
	RepoPath string
}

func RunReview(req ReviewRequest) (string, error) {
	return "review command invoked for repo: " + req.RepoPath, nil
}
