package review

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/bluekeyes/go-gitdiff/gitdiff"
	"github.com/mehditabet/exemplar-cli/internal/platform/git"
)

type Pipeline struct{}

type runState struct {
	request        ReviewRequest
	context        ReviewContext
	findings       []Finding
	executedStages []string
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

	state.context.RepositoryRoot = gitRoot

	p.recordStage(state, "resolve_repository")

	return nil
}

func (p Pipeline) resolveCurrentBranch(state *runState) error {

	branch, err := git.GetCurrentBranch(state.context.RepositoryRoot)

	if err != nil {
		return err
	}

	state.context.CurrentBranch = branch

	p.recordStage(state, "resolve_current_branch")
	return nil
}

func (p Pipeline) resolveChangedFiles(state *runState) error {

	changedFiles, err := git.GetChangedFiles(state.context.RepositoryRoot)

	if err != nil {
		return err
	}
	state.context.ChangedFiles = changedFiles

	p.recordStage(state, "resolve_changed_files")
	return nil
}

func (p Pipeline) resolveDiff(state *runState) error {

	diff, err := git.GetDiff(state.context.RepositoryRoot)

	if err != nil {
		return err
	}

	state.context.Diff = diff

	p.recordStage(state, "resolve_diff")
	return nil
}
func (p Pipeline) parseDiff(state *runState) error {

	if state.context.Diff == "" {
		state.context.ParseDiff = ParseDiff{Files: []DiffFile{}}
		p.recordStage(state, "parse_diff")
		return nil
	}
	files, _, err := gitdiff.Parse(strings.NewReader(state.context.Diff))

	if err != nil {
		return err
	}
	parseFiles := make([]DiffFile, 0, len(files))

	for _, file := range files {
		diffFile := DiffFile{
			OldPath:  file.OldName,
			NewPath:  file.NewName,
			IsNew:    file.IsNew,
			IsDelete: file.IsDelete,
			IsRename: file.IsRename,
			IsCopy:   file.IsCopy,
			Hunks:    make([]DiffHunk, 0, len(file.TextFragments)),
			isBinary: file.IsBinary,
		}

		for _, hunk := range file.TextFragments {
			hunk := DiffHunk{
				Header: hunk.Header(),
				Lines:  make([]DiffLine, 0, len(hunk.Lines)),
			}

			for _, line := range hunk.Lines {
				hunk.Lines = append(hunk.Lines, DiffLine{
					Operation: line.Operation,
					Content:   line.Content,
				})
			}
			diffFile.Hunks = append(diffFile.Hunks, hunk)
		}
		parseFiles = append(parseFiles, diffFile)
	}

	state.context.ParseDiff = ParseDiff{Files: parseFiles}
	p.recordStage(state, "parse_diff")
	return nil
}

func (p Pipeline) evaluateWorkspace(state *runState) error {

	if len(state.context.ChangedFiles) == 0 {
		state.context.Warnings = append(state.context.Warnings, "no changed files detected in the repository. The review will be based on the current state of the repository.")
	}
	if len(state.context.ChangedFiles) > 0 && state.context.Diff == "" {
		state.context.Warnings = append(state.context.Warnings, "changed files detected but no diff could be resolved. The review will be based on the current state of the repository.")
	}

	p.recordStage(state, "evaluate_workspace")

	return nil
}

func (p Pipeline) buildResult(ctx ReviewContext, executedStages []string) ReviewResult {
	stages := append(append([]string(nil), executedStages...), "build_result")

	return ReviewResult{
		RepositoryRoot: ctx.RepositoryRoot,
		Message:        "review command invoked for git repository: " + ctx.RepositoryRoot + " on branch: " + ctx.CurrentBranch + " with " + fmt.Sprint(len(ctx.ChangedFiles)) + " changed files.",
		Status:         "completed",
		CurrentBranch:  ctx.CurrentBranch,
		ChangedFiles:   ctx.ChangedFiles,
		ExecutedStages: stages,
		Diff:           ctx.Diff,
		Warnings:       ctx.Warnings,
	}
}

func (p Pipeline) BuildContext(req ReviewRequest) (ReviewContext, []string, error) {
	state := runState{request: req}

	err := p.validateRequest(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.resolveRepository(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.resolveCurrentBranch(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.resolveChangedFiles(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.resolveDiff(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.parseDiff(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	err = p.evaluateWorkspace(&state)

	if err != nil {
		return ReviewContext{}, nil, err
	}

	return state.context, append([]string(nil), state.executedStages...), nil
}

func (p Pipeline) Run(req ReviewRequest) (ReviewResult, error) {
	ctx, executedStages, err := p.BuildContext(req)
	if err != nil {
		return ReviewResult{}, err
	}

	result := p.buildResult(ctx, executedStages)
	return result, nil
}

func (p Pipeline) recordStage(state *runState, stageName string) {
	state.executedStages = append(state.executedStages, stageName)
}
