// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0

//go:build gui

package main

import (
	"fmt"
	"os"

	"github.com/micasa-dev/micasa/internal/data"
	"github.com/micasa-dev/micasa/internal/gui"
	"github.com/spf13/cobra"
)

const (
	rootShort = "A desktop GUI for tracking everything about your home"
	rootLong  = "A desktop GUI for tracking everything about your home."
)

// rootRun is the default action for the root command. With the gui build
// tag, this opens the Fyne desktop GUI.
func rootRun(_ *cobra.Command, opts *runOpts) error {
	return runGUI(opts)
}

// rootExtraSubcmds returns the tui subcommand for GUI builds so the terminal
// UI remains accessible via `micasa tui`.
func rootExtraSubcmds() []*cobra.Command {
	return []*cobra.Command{newTUICmd()}
}

// runGUI opens the Fyne desktop GUI with the resolved database.
func runGUI(opts *runOpts) error {
	dbPath, err := opts.resolveDBPath()
	if err != nil {
		return fmt.Errorf("resolve db path: %w", err)
	}
	if opts.printPath {
		_, _ = fmt.Fprintln(os.Stdout, dbPath)
		return nil
	}
	return launchGUI(dbPath)
}

// launchGUI opens the database, migrates, and starts the Fyne GUI.
func launchGUI(dbPath string) error {
	store, err := data.Open(dbPath)
	if err != nil {
		return fmt.Errorf("open database: %w", err)
	}
	if err := store.AutoMigrate(); err != nil {
		return fmt.Errorf("migrate database: %w", err)
	}
	if err := store.SeedDefaults(); err != nil {
		return fmt.Errorf("seed defaults: %w", err)
	}
	if err := store.ResolveCurrency(""); err != nil {
		return fmt.Errorf("resolve currency: %w", err)
	}
	gui.Run(store)
	return nil
}

// newTUICmd exposes the original terminal UI as the "tui" subcommand.
func newTUICmd() *cobra.Command {
	opts := &runOpts{}
	cmd := &cobra.Command{
		Use:           "tui [database-path]",
		Short:         "Launch the terminal UI (classic interface)",
		Long:          "Launch the original terminal UI for micasa.",
		Args:          cobra.MaximumNArgs(1),
		SilenceErrors: true,
		SilenceUsage:  true,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) > 0 {
				opts.dbPath = args[0]
			}
			return runTUI(cmd.OutOrStdout(), opts)
		},
	}
	cmd.Flags().
		BoolVar(&opts.printPath, "print-path", false, "Print the resolved database path and exit")
	return cmd
}
