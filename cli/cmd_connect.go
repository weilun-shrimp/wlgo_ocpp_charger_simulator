package cli

import "fmt"

func init() { register("connect", handleConnect) }

// handleConnect connects to the server and, on success, sends BootNotification
// followed by the current StatusNotification.
func handleConnect(ctx *CommandContext, args []string) {
	if ctx.Charger.IsConnected() {
		fmt.Fprintln(ctx.Out, "Already connected")
		return
	}
	if err := ctx.Charger.Connect(); err != nil {
		fmt.Fprintf(ctx.Out, "Error: %v\n", err)
		return
	}
	fmt.Fprintln(ctx.Out, "Connected to server")

	if err := ctx.Charger.BootNotification(); err != nil {
		fmt.Fprintf(ctx.Out, "BootNotification failed: %v\n", err)
		return
	}

	if err := ctx.Charger.StatusNotification(ctx.Charger.GetStatus()); err != nil {
		fmt.Fprintf(ctx.Out, "StatusNotification failed: %v\n", err)
	}
}
