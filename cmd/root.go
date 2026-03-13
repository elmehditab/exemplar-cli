package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// NewRootCmd creates the root command for the exemplar-cli application.
func NewRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "exemplar-cli",
		Short: "Smart CLI Reviews act as quality gates for Codex, Claude, Gemini, and you",
	}

	cmd.AddCommand(newReviewCmd())

	return cmd
}

// Execute runs the root command, handling any errors that occur.
func Execute() {
	err := NewRootCmd().Execute()
	if err != nil {
		os.Exit(1)
	}
}
