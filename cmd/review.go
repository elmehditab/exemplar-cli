package cmd

import (
	"github.com/mehditabet/exemplar-cli/internal/app"
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run a code review pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			message, err := app.RunReview()

			if err != nil {
				return err
			}

			cmd.Println(message)
			return nil
		},
	}
	return cmd
}
