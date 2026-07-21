package cli

import "fmt"

func init() {
	register("quit", handleQuit)
	register("exit", handleQuit)
}

// handleQuit handles both the "quit" and "exit" aliases. The interactive loop
// is intentionally exited via Ctrl+C, so this only prints guidance.
func handleQuit(ctx *CommandContext, args []string) {
	fmt.Fprintln(ctx.Out, "Use Ctrl+C to exit")
}
