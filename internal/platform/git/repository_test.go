package git

import (
	"os"
	"os/exec"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

// Tests that the repository root is resolved for a valid Git repository.
func TestResolveGitRepository(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)

	root, err := ResolveGitRepository(repoDir)
	if err != nil {
		t.Fatalf("ResolveGitRepository returned error: %v", err)
	}

	if resolvePath(t, root) != resolvePath(t, repoDir) {
		t.Fatalf("expected root %q, got %q", repoDir, root)
	}
}

// Tests that resolving the repository root fails for a non-Git directory.
func TestResolveGitRepositoryReturnsErrorForNonRepository(t *testing.T) {
	t.Parallel()

	_, err := ResolveGitRepository(t.TempDir())
	if err == nil {
		t.Fatal("expected error, got nil")
	}

	if !strings.Contains(err.Error(), "repo path is not a git repository") {
		t.Fatalf("unexpected error: %v", err)
	}
}

// Tests that the current branch is returned for a valid Git repository.
func TestGetCurrentBranch(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)

	branch, err := GetCurrentBranch(repoDir)
	if err != nil {
		t.Fatalf("GetCurrentBranch returned error: %v", err)
	}

	if branch != "main" {
		t.Fatalf("expected branch main, got %q", branch)
	}
}

// Tests that changed tracked and untracked files are both detected.
func TestGetChangedFiles(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)
	writeFile(t, repoDir, "tracked.txt", "first\n")
	gitRun(t, repoDir, "add", "tracked.txt")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	writeFile(t, repoDir, "tracked.txt", "second\n")
	writeFile(t, repoDir, "new.txt", "new\n")

	files, err := GetChangedFiles(repoDir)
	if err != nil {
		t.Fatalf("GetChangedFiles returned error: %v", err)
	}

	expected := []string{"new.txt", "tracked.txt"}
	if !reflect.DeepEqual(sortStrings(files), expected) {
		t.Fatalf("expected changed files %v, got %v", expected, files)
	}
}

// Tests that the diff is returned when a tracked file is modified.
func TestGetDiff(t *testing.T) {
	t.Parallel()

	repoDir := initGitRepository(t)
	writeFile(t, repoDir, "tracked.txt", "first\n")
	gitRun(t, repoDir, "add", "tracked.txt")
	gitRun(t, repoDir, "commit", "-m", "initial commit")

	writeFile(t, repoDir, "tracked.txt", "second\n")

	diff, err := GetDiff(repoDir)
	if err != nil {
		t.Fatalf("GetDiff returned error: %v", err)
	}

	if diff == "" {
		t.Fatal("expected non-empty diff")
	}

	if !strings.Contains(diff, "-first") || !strings.Contains(diff, "+second") {
		t.Fatalf("unexpected diff: %q", diff)
	}
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

func sortStrings(values []string) []string {
	tidy := append([]string(nil), values...)
	for i := range tidy {
		for j := i + 1; j < len(tidy); j++ {
			if tidy[j] < tidy[i] {
				tidy[i], tidy[j] = tidy[j], tidy[i]
			}
		}
	}

	return tidy
}
