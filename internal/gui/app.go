// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0
//go:build gui


// Package gui provides a Fyne-based desktop GUI for micasa.
package gui

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/micasa-dev/micasa/internal/data"
)

const windowTitle = "micasa"

// Run opens the main micasa GUI window using the provided data store.
// It blocks until the window is closed.
func Run(store *data.Store) {
	a := app.New()
	a.Settings().SetTheme(theme.DefaultTheme())

	w := a.NewWindow(windowTitle)
	w.Resize(fyne.NewSize(1100, 700))
	w.SetMaster()

	tabs := container.NewAppTabs(
		container.NewTabItemWithIcon("House", theme.HomeIcon(), newHouseTab(store, w)),
		container.NewTabItemWithIcon("Projects", theme.DocumentIcon(), newProjectsTab(store, w)),
		container.NewTabItemWithIcon("Vendors", theme.AccountIcon(), newVendorsTab(store, w)),
		container.NewTabItemWithIcon("Appliances", theme.ComputerIcon(), newAppliancesTab(store, w)),
		container.NewTabItemWithIcon("Maintenance", theme.ConfirmIcon(), newMaintenanceTab(store, w)),
		container.NewTabItemWithIcon("Incidents", theme.WarningIcon(), newIncidentsTab(store, w)),
		container.NewTabItemWithIcon("Service Logs", theme.ListIcon(), newServiceLogsTab(store, w)),
	)
	tabs.SetTabLocation(container.TabLocationLeading)

	status := widget.NewLabel("Ready")
	status.Alignment = fyne.TextAlignLeading

	w.SetContent(container.NewBorder(nil, status, nil, nil, tabs))
	w.ShowAndRun()
}
