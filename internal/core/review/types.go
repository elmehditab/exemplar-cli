package review

type ReviewRequest struct {
	RepoPath string
}

type ReviewResult struct {
	RepositoryRoot string
	Message        string
	Status         string
	ExecutedStages []string
	Warnings       []string
}
