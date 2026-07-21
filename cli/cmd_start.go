package cli

import "fmt"

func init() { register("start", handleStart) }

// handleStart starts a transaction for the given idTag.
func handleStart(ctx *CommandContext, args []string) {
	if len(args) < 1 {
		fmt.Fprintln(ctx.Out, "Usage: start <idTag>")
		return
	}
	idTag := args[0]
	if err := ctx.Charger.StartTransaction(idTag); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
	} else {
		fmt.Fprintln(ctx.Out, "Transaction started")
	}
}
