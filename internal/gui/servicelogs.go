// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0
//go:build gui


package gui

import (
	"fmt"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/micasa-dev/micasa/internal/data"
	"github.com/micasa-dev/micasa/internal/uid"
)

func newServiceLogsTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var entries []data.ServiceLogEntry
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(entries) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Maintenance Item")
				case 1:
					lbl.SetText("Serviced At")
				case 2:
					lbl.SetText("Vendor")
				case 3:
					lbl.SetText("Cost")
				case 4:
					lbl.SetText("Notes")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(entries) {
				lbl.SetText("")
				return
			}
			e := entries[i]
			switch id.Col {
			case 0:
				lbl.SetText(e.MaintenanceItem.Name)
			case 1:
				lbl.SetText(e.ServicedAt.Format("2006-01-02"))
			case 2:
				lbl.SetText(e.Vendor.Name)
			case 3:
				lbl.SetText(centsStr(e.CostCents))
			case 4:
				lbl.SetText(e.Notes)
			}
		},
	)
	table.SetColumnWidth(0, 250)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 180)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 250)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		entries, err = store.ListAllServiceLogEntries(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showServiceLogDialog(store, w, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(entries) {
			dialog.ShowInformation("No selection", "Please select a service log entry to delete.", w)
			return
		}
		e := entries[selected]
		dialog.ShowConfirm("Delete Service Log Entry",
			fmt.Sprintf("Delete service log entry for %q on %s?",
				e.MaintenanceItem.Name,
				e.ServicedAt.Format("2006-01-02"),
			),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteServiceLog(e.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showServiceLogDialog(store *data.Store, w fyne.Window, onSave func()) {
	items, err := store.ListMaintenance(false)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	vendors, err := store.ListVendors(false)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	itemNames := make([]string, len(items))
	itemIDs := make([]string, len(items))
	for i, m := range items {
		itemNames[i] = m.Name
		itemIDs[i] = m.ID
	}

	vendorNames := make([]string, len(vendors)+1)
	vendorObjs := make([]data.Vendor, len(vendors)+1)
	vendorNames[0] = "(none)"
	for i, v := range vendors {
		vendorNames[i+1] = v.Name
		vendorObjs[i+1] = v
	}

	itemSelect := widget.NewSelect(itemNames, nil)
	servicedAt := widget.NewEntry()
	servicedAt.SetText(time.Now().Format("2006-01-02"))
	vendorSelect := widget.NewSelect(vendorNames, nil)
	vendorSelect.SetSelectedIndex(0)
	cost := widget.NewEntry()
	cost.SetPlaceHolder("e.g. 100.00")
	notes := widget.NewMultiLineEntry()
	notes.SetMinRowsVisible(3)

	if len(itemNames) > 0 {
		itemSelect.SetSelectedIndex(0)
	}

	form := widget.NewForm(
		widget.NewFormItem("Maintenance Item *", itemSelect),
		widget.NewFormItem("Serviced At", servicedAt),
		widget.NewFormItem("Vendor", vendorSelect),
		widget.NewFormItem("Cost ($)", cost),
		widget.NewFormItem("Notes", notes),
	)

	d := dialog.NewCustomConfirm("Add Service Log Entry", "Save", "Cancel", form,
		func(confirmed bool) {
			if !confirmed {
				return
			}
			itemIdx := itemSelect.SelectedIndex()
			if itemIdx < 0 {
				dialog.ShowInformation("Validation", "Maintenance item is required.", w)
				return
			}

			svcDate := time.Now()
			if t := parseDateOr(servicedAt.Text); t != nil {
				svcDate = *t
			}

			vendorIdx := vendorSelect.SelectedIndex()
			vendor := vendorObjs[0] // empty vendor = no vendor
			if vendorIdx > 0 {
				vendor = vendorObjs[vendorIdx]
			}

			var costCents *int64
			if cost.Text != "" {
				v := parseFloatOr(cost.Text, -1)
				if v >= 0 {
					c := int64(v * 100)
					costCents = &c
				}
			}

			now := time.Now()
			entry := data.ServiceLogEntry{
				ID:                uid.New(),
				MaintenanceItemID: itemIDs[itemIdx],
				ServicedAt:        svcDate,
				CostCents:         costCents,
				Notes:             notes.Text,
				CreatedAt:         now,
				UpdatedAt:         now,
			}
			if err := store.CreateServiceLog(&entry, vendor); err != nil {
				dialog.ShowError(err, w)
				return
			}
			if onSave != nil {
				onSave()
			}
		}, w)
	d.Resize(fyne.NewSize(500, 380))
	d.Show()
}
