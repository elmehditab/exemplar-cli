package cmd

import (
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

			cmd.Println("Repository Root:", result.RepositoryRoot)
			cmd.Println("Current Branch:", result.CurrentBranch)
			cmd.Println("Changed Files:", result.ChangedFiles)
			cmd.Println("Diff:", result.Diff)
			cmd.Println("Warnings:", result.Warnings)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Path to the repository")
	return cmd
}
