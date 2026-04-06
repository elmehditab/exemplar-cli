package review

type ReviewRequest struct {
	RepoPath string
}

type ReviewContext struct {
	RepositoryRoot string
	CurrentBranch  string
	ChangedFiles   []string
	Diff           string
	Warnings       []string
}

type ReviewResult struct {
	RepositoryRoot string
	Message        string
	Status         string
	CurrentBranch  string
	ChangedFiles   []string
	ExecutedStages []string
	Diff           string
	Warnings       []string
}
