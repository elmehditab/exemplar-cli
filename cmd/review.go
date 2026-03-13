package cmd

import (
	"github.com/spf13/cobra"
)

func newReviewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "review",
		Short: "Run a code review pipeline",
		RunE: func(cmd *cobra.Command, args []string) error {
			cmd.Println("review command invoked")
			return nil
		},
	}

	return cmd
}
