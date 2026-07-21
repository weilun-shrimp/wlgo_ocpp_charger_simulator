package cli

import "fmt"

func init() { register("stop", handleStop) }

// handleStop stops the current transaction, defaulting the reason to "Local".
func handleStop(ctx *CommandContext, args []string) {
	reason := "Local"
	if len(args) >= 1 {
		reason = args[0]
	}
	if err := ctx.Charger.StopTransaction(reason); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintln(ctx.Out, "Transaction stopped")
	}
}
