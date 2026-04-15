// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0

//go:build !gui

package main

import (
"github.com/spf13/cobra"
)

const (
rootShort = "A terminal UI for tracking everything about your home"
rootLong  = "A terminal UI for tracking everything about your home."
)

// rootRun is the default action for the root command. Without the gui build
// tag, this launches the TUI.
func rootRun(cmd *cobra.Command, opts *runOpts) error {
return runTUI(cmd.OutOrStdout(), opts)
}

// rootExtraSubcmds returns no additional subcommands in the TUI-only build.
func rootExtraSubcmds() []*cobra.Command { return nil }
