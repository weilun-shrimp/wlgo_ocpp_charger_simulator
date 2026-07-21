package cli

import (
	"io"

	"github.com/weilun-shrimp/wlgo_ocpp_charger_simulator/config"
)

// CommandContext carries everything a command handler is allowed to touch: the
// charger it acts on, the configuration (for version- and limit-dependent
// output), and the writer all output must go to. Handlers write to Out rather
// than os.Stdout directly so tests can capture and assert their output.
type CommandContext struct {
	Charger Charger
	Config  *config.Config
	Out     io.Writer
}
