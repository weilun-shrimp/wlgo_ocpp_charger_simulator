// Package cli implements the interactive command layer for the OCPP charger
// simulator. It parses lines from an input stream, dispatches each command to a
// dedicated handler registered in a command table, and writes all output to an
// injected writer so that every command is independently testable.
package cli
