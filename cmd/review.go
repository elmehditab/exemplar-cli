package cmd

import (
	"fmt"

	"github.com/mehditabet/exemplar-cli/internal/app"
	"github.com/mehditabet/exemplar-cli/internal/core/review"

	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {

	var repoPath string

	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run a code review pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {

			req := review.ReviewRequest{
				RepoPath: repoPath,
			}

			result, err := app.RunReview(req)

			if err != nil {
				return err
			}

			printReviewResult(cmd, result)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Path to the repository")
	return cmd
}

func printReviewResult(cmd *cobra.Command, result review.ReviewResult) {
	cmd.Println("Review Summary")
	cmd.Println("Status:", result.Status)
	cmd.Println("Repository:", result.RepositoryRoot)
	cmd.Println("Branch:", result.CurrentBranch)
	cmd.Println("Changed files:", len(result.ChangedFiles))
	cmd.Println("Parsed files:", result.ParsedDiff.Stats.FilesChanged)
	cmd.Println("Review targets:", len(result.ParsedDiff.ReviewTargets))
	cmd.Println("Lines added:", result.ParsedDiff.Stats.LinesAdded)
	cmd.Println("Lines deleted:", result.ParsedDiff.Stats.LinesDeleted)
	cmd.Println("Binary files:", result.ParsedDiff.Stats.BinaryFiles)

	if len(result.ChangedFiles) > 0 {
		cmd.Println()
		cmd.Println("Changed Files")
		for _, file := range result.ChangedFiles {
			cmd.Println("-", file)
		}
	}

	if len(result.ParsedDiff.Files) > 0 {
		cmd.Println()
		cmd.Println("Parsed Diff")
		for _, file := range result.ParsedDiff.Files {
			cmd.Println("-", file.Status, displayDiffPath(file), "("+lineStats(file)+")")
		}
	}

	if len(result.Warnings) > 0 {
		cmd.Println()
		cmd.Println("Warnings")
		for _, warning := range result.Warnings {
			cmd.Println("-", warning)
		}
	}
}

func displayDiffPath(file review.DiffFile) string {
	if file.NewPath != "" {
		return file.NewPath
	}
	return file.OldPath
}

func lineStats(file review.DiffFile) string {
	return fmt.Sprintf("+%d -%d", file.LinesAdded, file.LinesDeleted)
}
