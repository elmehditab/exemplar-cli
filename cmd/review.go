package cmd

import (
	"github.com/mehditabet/exemplar-cli/internal/app"

	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {

	var repoPath string

	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run a code review pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {

			req := app.ReviewRequest{
				RepoPath: repoPath,
			}

			message, err := app.RunReview(req)

			if err != nil {
				return err
			}

			cmd.Println(message)
			return nil
		},
	}

	cmd.Flags().StringVar(&repoPath, "repo", ".", "Path to the repository")
	return cmd
}
