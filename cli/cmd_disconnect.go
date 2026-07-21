package cli

import "fmt"

func init() { register("disconnect", handleDisconnect) }

// handleDisconnect disconnects from the server if currently connected.
func handleDisconnect(ctx *CommandContext, args []string) {
	if !ctx.Charger.IsConnected() {
		fmt.Fprintln(ctx.Out, "Not connected")
		return
	}
	ctx.Charger.Disconnect()
	fmt.Fprintln(ctx.Out, "Disconnected from server")
}
