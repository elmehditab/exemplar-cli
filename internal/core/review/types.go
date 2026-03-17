package review

type ReviewRequest struct {
	RepoPath string
}

type ReviewResult struct {
	RepositoryRoot string
	Message        string
	Status         string
	CurrentBranch  string
	ExecutedStages []string
	Warnings       []string
}
