package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

const (
	colorReset = "\033[0m"
	colorDim   = "\033[2m"
	colorCyan  = "\033[36m"
	colorBlue  = "\033[34m"
	colorGreen = "\033[32m"
	colorGold  = "\033[33m"
	colorWhite = "\033[97m"
)

func printWelcome(cmd *cobra.Command) {
	out := cmd.OutOrStdout()

	fmt.Fprint(out, colorCyan)
	fmt.Fprintln(out, "Welcome to Exemplar CLI "+colorWhite+"v"+cmd.Root().Version+colorReset)
	fmt.Fprint(out, colorDim)
	fmt.Fprintln(out, "----------------------------------------------------------------------------")
	fmt.Fprint(out, colorReset)
	printLogo(out)
	fmt.Fprint(out, colorDim)
	fmt.Fprintln(out, "----------------------------------------------------------------------------")
	fmt.Fprint(out, colorReset)
	fmt.Fprintln(out, colorWhite+"Review code with repository context before generating findings."+colorReset)
	fmt.Fprintln(out)
	fmt.Fprintln(out, colorGold+"Pipeline"+colorReset+"  collect changes -> parse diff -> build review targets -> report")
	fmt.Fprintln(out, colorGreen+"Command "+colorReset+"  exemplar-cli review --repo .")
	fmt.Fprintln(out, colorDim+"Status  "+colorReset+"  MVP parser ready for evidence-backed findings")
}

func printLogo(out io.Writer) {
	fmt.Fprintln(out)
	fmt.Fprintln(out, colorBlue+"        ______                         __"+colorReset)
	fmt.Fprintln(out, colorBlue+"       / ____/  _____  ____ ___  ____  / /___ ______"+colorReset)
	fmt.Fprintln(out, colorBlue+"      / __/ | |/_/ _ \\/ __ `__ \\/ __ \\/ / __ `/ ___/"+colorReset)
	fmt.Fprintln(out, colorCyan+"     / /____>  </  __/ / / / / / /_/ / / /_/ / /"+colorReset)
	fmt.Fprintln(out, colorCyan+"    /_____/_/|_|\\___/_/ /_/ /_/ .___/_/\\__,_/_/"+colorReset)
	fmt.Fprintln(out, colorCyan+"                              /_/              "+colorReset)
	fmt.Fprintln(out)
	fmt.Fprintln(out, colorGold+"       *       "+colorDim+"repo-aware code review for AI-generated changes"+colorReset)
	fmt.Fprintln(out, colorDim+"              parse diffs, anchor evidence, and review with context"+colorReset)
	fmt.Fprintln(out)
}
