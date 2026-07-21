package cli

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

// Handler executes a single interactive command. args holds the
// whitespace-split tokens that follow the command word (the verb itself is
// excluded). All user-facing output must be written to ctx.Out.
type Handler func(ctx *CommandContext, args []string)

// registry maps a lower-cased command verb to its handler. Command files
// populate it from their init() functions via register().
var registry = map[string]Handler{}

// register adds a handler for the given command verb. It panics on a duplicate
// registration so that a wiring mistake fails loudly at startup rather than
// silently shadowing a command.
func register(name string, h Handler) {
	if _, exists := registry[name]; exists {
		panic("cli: duplicate command registration: " + name)
	}
	registry[name] = h
}

// parseCommand splits a raw input line into a lower-cased command verb and its
// arguments. ok is false when the line is empty or whitespace-only, in which
// case no command should run.
func parseCommand(line string) (cmd string, args []string, ok bool) {
	fields := strings.Fields(strings.TrimSpace(line))
	if len(fields) == 0 {
		return "", nil, false
	}
	return strings.ToLower(fields[0]), fields[1:], true
}

// Dispatch runs the handler registered for cmd. When no handler is registered
// it prints the unknown-command message to ctx.Out.
func Dispatch(ctx *CommandContext, cmd string, args []string) {
	h, ok := registry[cmd]
	if !ok {
		fmt.Fprintf(ctx.Out, "Unknown command: %s. Type 'help' for available commands.\n", cmd)
		return
	}
	h(ctx, args)
}

// Run reads command lines from in, printing a prompt to out before each read,
// and dispatches each parsed command. It returns when in is exhausted (EOF),
// which lets tests drive a finite sequence of commands and then assert results.
func Run(c Charger, cfg *config.Config, in io.Reader, out io.Writer) {
	ctx := &CommandContext{Charger: c, Config: cfg, Out: out}
	reader := bufio.NewReader(in)

	for {
		fmt.Fprint(out, "> ")
		line, err := reader.ReadString('\n')
		if line != "" {
			if cmd, args, ok := parseCommand(line); ok {
				Dispatch(ctx, cmd, args)
			}
		}
		if err != nil {
			return
		}
	}
}
