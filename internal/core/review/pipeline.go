package review

import (
	"errors"
	"fmt"
	"os"

	"github.com/mehditabet/exemplar-cli/internal/platform/git"
)

type Pipeline struct{}

type runState struct {
	request        ReviewRequest
	repositoryRoot string
	currentBranch  string
	changedFiles   []string
	executedStages []string
	warnings       []string
	diff           string
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

	p.recordStage(state, "validate_request")

	return nil
}

func (p Pipeline) resolveRepository(state *runState) error {

	gitRoot, err := git.ResolveGitRepository(state.request.RepoPath)

	if err != nil {
		return err
	}

	state.repositoryRoot = gitRoot

	p.recordStage(state, "resolve_repository")

	return nil
}

func (p Pipeline) resolveCurrentBranch(state *runState) error {

	branch, err := git.GetCurrentBranch(state.repositoryRoot)

	if err != nil {
		return err
	}

	state.currentBranch = branch

	p.recordStage(state, "resolve_current_branch")
	return nil
}

func (p Pipeline) resolveChangedFiles(state *runState) error {

	changedFiles, err := git.GetChangedFiles(state.repositoryRoot)

	if err != nil {
		return err
	}
	state.changedFiles = changedFiles

	p.recordStage(state, "resolve_changed_files")
	return nil
}

func (p Pipeline) resolveDiff(state *runState) error {

	diff, err := git.GetDiff(state.repositoryRoot)

	if err != nil {
		return err
	}

	state.diff = diff

	p.recordStage(state, "resolve_diff")
	return nil
}

func (p Pipeline) evaluateWorkspace(state *runState) error {

	if len(state.changedFiles) == 0 {
		state.warnings = append(state.warnings, "no changed files detected in the repository. The review will be based on the current state of the repository.")
	}
	if len(state.changedFiles) > 0 && state.diff == "" {
		state.warnings = append(state.warnings, "changed files detected but no diff could be resolved. The review will be based on the current state of the repository.")
	}

	p.recordStage(state, "evaluate_workspace")

	return nil
}

func (p Pipeline) buildResult(state *runState) ReviewResult {

	p.recordStage(state, "build_result")

	return ReviewResult{
		RepositoryRoot: state.repositoryRoot,
		Message:        "review command invoked for git repository: " + state.repositoryRoot + " on branch: " + state.currentBranch + " with " + fmt.Sprint(len(state.changedFiles)) + " changed files.",
		Status:         "completed",
		CurrentBranch:  state.currentBranch,
		ExecutedStages: state.executedStages,
		ChangedFiles:   state.changedFiles,
		Warnings:       state.warnings,
		Diff:           state.diff,
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

	err = p.resolveCurrentBranch(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	err = p.resolveChangedFiles(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	err = p.resolveDiff(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	err = p.evaluateWorkspace(&state)

	if err != nil {
		return ReviewResult{}, err
	}

	result := p.buildResult(&state)

	return result, nil
}

func (p Pipeline) recordStage(state *runState, stageName string) {
	state.executedStages = append(state.executedStages, stageName)
}
