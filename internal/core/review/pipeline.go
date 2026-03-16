package review

import (
	"errors"
	"os"

	"github.com/mehditabet/exemplar-cli/internal/platform/git"
)

type Pipeline struct{}

type runState struct {
	request        ReviewRequest
	repositoryRoot string
}

func (p Pipeline) validateRequest(state *runState) error {

	if state.request.RepoPath == "" {
		return errors.New("repo path is required")
	}

	info, err := os.Stat(state.request.RepoPath)

	if err != nil {
		return errors.New("repo path does not exist: " + state.request.RepoPath)
	}

	if !info.IsDir() {
		return errors.New("repo path must be a directory")
	}

	return nil
}

func (p Pipeline) resolveRepository(state *runState) error {

	gitRoot, err := git.ResolveGitRepository(state.request.RepoPath)

	if err != nil {
		return err
	}

	state.repositoryRoot = gitRoot

	return nil
}

func (p Pipeline) buildResult(state runState) ReviewResult {

	return ReviewResult{
		RepositoryRoot: state.repositoryRoot,
		Message:        "review command invoked for git repository: " + state.repositoryRoot,
	}
}

func (p Pipeline) Run(req ReviewRequest) (ReviewResult, error) {

	state := runState{request: req}

	err := p.validateRequest(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	err = p.resolveRepository(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	result := p.buildResult(state)

	return result, nil
}
