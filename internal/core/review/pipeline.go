package review

import (
	"errors"
	"os"

	"github.com/mehditabet/exemplar-cli/internal/platform/git"
)

type Pipeline struct{}

func (p Pipeline) Run(req ReviewRequest) (ReviewResult, error) {

	if req.RepoPath == "" {
		return ReviewResult{}, errors.New("repo path is required")
	}

	info, err := os.Stat(req.RepoPath)
	if err != nil {
		return ReviewResult{}, errors.New("repo path does not exist: " + req.RepoPath)
	}

	if !info.IsDir() {
		return ReviewResult{}, errors.New("repo path must be a directory")
	}

	gitRoot, err := git.ResolveRepositoryRoot(req.RepoPath)
	if err != nil {
		return ReviewResult{}, err
	}

	return ReviewResult{
		RepositoryRoot: gitRoot,
		Message:        "review command invoked for git repository: " + gitRoot,
	}, nil
}
