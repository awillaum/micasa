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

func newMaintenanceTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var items []data.MaintenanceItem
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(items) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Name")
				case 1:
					lbl.SetText("Category")
				case 2:
					lbl.SetText("Season")
				case 3:
					lbl.SetText("Interval (mo)")
				case 4:
					lbl.SetText("Due Date")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(items) {
				lbl.SetText("")
				return
			}
			m := items[i]
			switch id.Col {
			case 0:
				lbl.SetText(m.Name)
			case 1:
				lbl.SetText(m.Category.Name)
			case 2:
				lbl.SetText(m.Season)
			case 3:
				lbl.SetText(fmt.Sprintf("%d", m.IntervalMonths))
			case 4:
				lbl.SetText(timeStr(m.DueDate))
			}
		},
	)
	table.SetColumnWidth(0, 250)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 110)
	table.SetColumnWidth(4, 110)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		items, err = store.ListMaintenance(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showMaintenanceDialog(store, w, nil, reload)
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(items) {
			dialog.ShowInformation("No selection", "Please select a maintenance item to edit.", w)
			return
		}
		item := items[selected]
		showMaintenanceDialog(store, w, &item, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(items) {
			dialog.ShowInformation("No selection", "Please select a maintenance item to delete.", w)
			return
		}
		item := items[selected]
		dialog.ShowConfirm("Delete Maintenance Item",
			fmt.Sprintf("Delete maintenance item %q?", item.Name),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteMaintenance(item.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showMaintenanceDialog(
	store *data.Store,
	w fyne.Window,
	existing *data.MaintenanceItem,
	onSave func(),
) {
	categories, err := store.MaintenanceCategories()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	catNames := make([]string, len(categories))
	catIDs := make([]string, len(categories))
	for i, c := range categories {
		catNames[i] = c.Name
		catIDs[i] = c.ID
	}

	appliances, err := store.ListAppliances(false)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	appNames := make([]string, len(appliances)+1)
	appIDs := make([]string, len(appliances)+1)
	appNames[0] = "(none)"
	appIDs[0] = ""
	for i, a := range appliances {
		appNames[i+1] = a.Name
		appIDs[i+1] = a.ID
	}

	name := widget.NewEntry()
	catSelect := widget.NewSelect(catNames, nil)
	appSelect := widget.NewSelect(appNames, nil)
	seasonSelect := widget.NewSelect(seasons(), nil)
	interval := widget.NewEntry()
	interval.SetPlaceHolder("months")
	dueDate := widget.NewEntry()
	dueDate.SetPlaceHolder("YYYY-MM-DD")
	manualURL := widget.NewEntry()
	notes := widget.NewMultiLineEntry()
	notes.SetMinRowsVisible(3)

	if existing != nil {
		name.SetText(existing.Name)
		for i, id := range catIDs {
			if id == existing.CategoryID {
				catSelect.SetSelectedIndex(i)
				break
			}
		}
		appSelect.SetSelectedIndex(0)
		if existing.ApplianceID != nil {
			for i, id := range appIDs {
				if id == *existing.ApplianceID {
					appSelect.SetSelectedIndex(i)
					break
				}
			}
		}
		seasonSelect.SetSelected(existing.Season)
		interval.SetText(fmt.Sprintf("%d", existing.IntervalMonths))
		dueDate.SetText(timeStr(existing.DueDate))
		manualURL.SetText(existing.ManualURL)
		notes.SetText(existing.Notes)
	} else {
		if len(catNames) > 0 {
			catSelect.SetSelectedIndex(0)
		}
		appSelect.SetSelectedIndex(0)
		seasonSelect.SetSelected(data.SeasonSpring)
		interval.SetText("12")
	}

	form := widget.NewForm(
		widget.NewFormItem("Name *", name),
		widget.NewFormItem("Category *", catSelect),
		widget.NewFormItem("Appliance", appSelect),
		widget.NewFormItem("Season", seasonSelect),
		widget.NewFormItem("Interval (months)", interval),
		widget.NewFormItem("Due Date", dueDate),
		widget.NewFormItem("Manual URL", manualURL),
		widget.NewFormItem("Notes", notes),
	)

	label := "Add Maintenance Item"
	if existing != nil {
		label = "Edit Maintenance Item"
	}

	d := dialog.NewCustomConfirm(label, "Save", "Cancel", form,
		func(confirmed bool) {
			if !confirmed {
				return
			}
			if name.Text == "" {
				dialog.ShowInformation("Validation", "Name is required.", w)
				return
			}
			catIdx := catSelect.SelectedIndex()
			if catIdx < 0 {
				dialog.ShowInformation("Validation", "Category is required.", w)
				return
			}

			intervalMonths := parseIntOr(interval.Text, 12)
			dueT := parseDateOr(dueDate.Text)

			appIdx := appSelect.SelectedIndex()
			var appID *string
			if appIdx > 0 {
				id := appIDs[appIdx]
				appID = &id
			}

			if existing != nil {
				existing.Name = name.Text
				existing.CategoryID = catIDs[catIdx]
				existing.ApplianceID = appID
				existing.Season = seasonSelect.Selected
				existing.IntervalMonths = intervalMonths
				existing.DueDate = dueT
				existing.ManualURL = manualURL.Text
				existing.Notes = notes.Text
				if err := store.UpdateMaintenance(*existing); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				now := time.Now()
				item := data.MaintenanceItem{
					ID:             uid.New(),
					Name:           name.Text,
					CategoryID:     catIDs[catIdx],
					ApplianceID:    appID,
					Season:         seasonSelect.Selected,
					IntervalMonths: intervalMonths,
					DueDate:        dueT,
					ManualURL:      manualURL.Text,
					Notes:          notes.Text,
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				if err := store.CreateMaintenance(&item); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}
			if onSave != nil {
				onSave()
			}
		}, w)
	d.Resize(fyne.NewSize(500, 500))
	d.Show()
}

func seasons() []string {
	return []string{
		data.SeasonSpring,
		data.SeasonSummer,
		data.SeasonFall,
		data.SeasonWinter,
	}
}
