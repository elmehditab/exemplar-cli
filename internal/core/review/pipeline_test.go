package review

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Tests that the pipeline returns the expected errors for invalid input.
func TestPipelineRunValidationErrors(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name    string
		req     ReviewRequest
		wantErr string
	}{
		{
			name:    "missing repo path",
			req:     ReviewRequest{},
			wantErr: "repo path is required",
		},
		{
			name:    "repo path does not exist",
			req:     ReviewRequest{RepoPath: filepath.Join(t.TempDir(), "missing")},
			wantErr: "repo path does not exist",
		},
		{
			name:    "repo path is not a directory",
			req:     ReviewRequest{RepoPath: createTempFile(t)},
			wantErr: "repo path must be a directory",
		},
		{
			name:    "repo path is not a git repository",
			req:     ReviewRequest{RepoPath: t.TempDir()},
			wantErr: "repo path is not a git repository",
		},
	}

	for _, tt := range tests {
		tt := tt
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			pipeline := Pipeline{}

			_, err := pipeline.Run(tt.req)
			if err == nil {
				t.Fatalf("expected error %q, got nil", tt.wantErr)
			}

			if !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected error containing %q, got %q", tt.wantErr, err.Error())
			}
		})
	}
}

// Tests that the pipeline returns a completed result for a clean repository.
func TestPipelineRunCleanRepository(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)
	writeFile(t, repoDir, "README.md", "# exemplar\n")
	gitRun(t, repoDir, "add", "README.md")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	result, err := Pipeline{}.Run(ReviewRequest{RepoPath: repoDir})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if resolvePath(t, result.RepositoryRoot) != resolvePath(t, repoDir) {
		t.Fatalf("expected repository root %q, got %q", repoDir, result.RepositoryRoot)
	}

	if result.CurrentBranch != "main" {
		t.Fatalf("expected current branch main, got %q", result.CurrentBranch)
	}

	if result.Status != "completed" {
		t.Fatalf("expected status completed, got %q", result.Status)
	}

	if len(result.ChangedFiles) != 0 {
		t.Fatalf("expected no changed files, got %v", result.ChangedFiles)
	}

	if result.Diff != "" {
		t.Fatalf("expected empty diff, got %q", result.Diff)
	}

	expectedStages := []string{
		"validate_request",
		"resolve_repository",
		"resolve_current_branch",
		"resolve_changed_files",
		"resolve_diff",
		"evaluate_workspace",
		"build_result",
	}
	if !reflect.DeepEqual(result.ExecutedStages, expectedStages) {
		t.Fatalf("expected stages %v, got %v", expectedStages, result.ExecutedStages)
	}

	expectedWarnings := []string{
		"no changed files detected in the repository. The review will be based on the current state of the repository.",
	}
	if !reflect.DeepEqual(result.Warnings, expectedWarnings) {
		t.Fatalf("expected warnings %v, got %v", expectedWarnings, result.Warnings)
	}

	if !strings.Contains(result.Message, repoDir) {
		t.Fatalf("expected message to mention repo root %q, got %q", repoDir, result.Message)
	}
}

// Tests that the pipeline detects changed files and returns the diff.
func TestPipelineRunRepositoryWithChanges(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)
	writeFile(t, repoDir, "main.go", "package main\n\nfunc main() {}\n")
	gitRun(t, repoDir, "add", "main.go")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	writeFile(t, repoDir, "main.go", "package main\n\nfunc main() {\n\tprintln(\"changed\")\n}\n")

	result, err := Pipeline{}.Run(ReviewRequest{RepoPath: repoDir})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}

	if len(result.ChangedFiles) != 1 {
		t.Fatalf("expected 1 changed file, got %v", result.ChangedFiles)
	}

	if result.ChangedFiles[0] != "main.go" {
		t.Fatalf("expected changed file main.go, got %v", result.ChangedFiles)
	}

	if result.Diff == "" {
		t.Fatal("expected non-empty diff")
	}

	if !strings.Contains(result.Diff, `+	println("changed")`) {
		t.Fatalf("expected diff to contain changed line, got %q", result.Diff)
	}

	if len(result.Warnings) != 0 {
		t.Fatalf("expected no warnings, got %v", result.Warnings)
	}
}

func createTempFile(t *testing.T) string {
	t.Helper()

	file, err := os.CreateTemp(t.TempDir(), "review-test-*")
	if err != nil {
		t.Fatalf("CreateTemp returned error: %v", err)
	}

	if err := file.Close(); err != nil {
		t.Fatalf("Close returned error: %v", err)
	}

	return file.Name()
}

func initGitRepository(t *testing.T) string {
	t.Helper()

	repoDir := t.TempDir()
	gitRun(t, repoDir, "init", "-b", "main")
	gitRun(t, repoDir, "config", "user.name", "Exemplar Tests")
	gitRun(t, repoDir, "config", "user.email", "tests@example.com")

	return repoDir
}

func writeFile(t *testing.T, repoDir, relativePath, content string) {
	t.Helper()

	fullPath := filepath.Join(repoDir, relativePath)
	if err := os.MkdirAll(filepath.Dir(fullPath), 0o755); err != nil {
		t.Fatalf("MkdirAll returned error: %v", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
}

func gitRun(t *testing.T, repoDir string, args ...string) string {
	t.Helper()

	cmd := exec.Command("git", append([]string{"-C", repoDir}, args...)...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("git %v returned error: %v\noutput: %s", args, err, string(out))
	}

	return strings.TrimSpace(string(out))
}

func resolvePath(t *testing.T, path string) string {
	t.Helper()

	resolved, err := filepath.EvalSymlinks(path)
	if err != nil {
		t.Fatalf("EvalSymlinks returned error for %q: %v", path, err)
	}

	return resolved
}
