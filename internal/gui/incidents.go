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

func newIncidentsTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var incidents []data.Incident
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(incidents) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Title")
				case 1:
					lbl.SetText("Status")
				case 2:
					lbl.SetText("Severity")
				case 3:
					lbl.SetText("Date Noticed")
				case 4:
					lbl.SetText("Location")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(incidents) {
				lbl.SetText("")
				return
			}
			inc := incidents[i]
			switch id.Col {
			case 0:
				lbl.SetText(inc.Title)
			case 1:
				lbl.SetText(inc.Status)
			case 2:
				lbl.SetText(inc.Severity)
			case 3:
				lbl.SetText(inc.DateNoticed.Format("2006-01-02"))
			case 4:
				lbl.SetText(inc.Location)
			}
		},
	)
	table.SetColumnWidth(0, 250)
	table.SetColumnWidth(1, 100)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 120)
	table.SetColumnWidth(4, 150)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		incidents, err = store.ListIncidents(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showIncidentDialog(store, w, nil, reload)
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(incidents) {
			dialog.ShowInformation("No selection", "Please select an incident to edit.", w)
			return
		}
		inc := incidents[selected]
		showIncidentDialog(store, w, &inc, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(incidents) {
			dialog.ShowInformation("No selection", "Please select an incident to delete.", w)
			return
		}
		inc := incidents[selected]
		dialog.ShowConfirm("Delete Incident",
			fmt.Sprintf("Delete incident %q?", inc.Title),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteIncident(inc.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showIncidentDialog(
	store *data.Store,
	w fyne.Window,
	existing *data.Incident,
	onSave func(),
) {
	appliances, err := store.ListAppliances(false)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	vendors, err := store.ListVendors(false)
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

	vendorNames := make([]string, len(vendors)+1)
	vendorIDs := make([]string, len(vendors)+1)
	vendorNames[0] = "(none)"
	vendorIDs[0] = ""
	for i, v := range vendors {
		vendorNames[i+1] = v.Name
		vendorIDs[i+1] = v.ID
	}

	title := widget.NewEntry()
	statusSelect := widget.NewSelect(incidentStatuses(), nil)
	severitySelect := widget.NewSelect(incidentSeverities(), nil)
	dateNoticed := widget.NewEntry()
	dateNoticed.SetPlaceHolder("YYYY-MM-DD")
	location := widget.NewEntry()
	appSelect := widget.NewSelect(appNames, nil)
	vendorSelect := widget.NewSelect(vendorNames, nil)
	cost := widget.NewEntry()
	cost.SetPlaceHolder("e.g. 250.00")
	description := widget.NewMultiLineEntry()
	description.SetMinRowsVisible(3)
	notes := widget.NewMultiLineEntry()
	notes.SetMinRowsVisible(2)

	if existing != nil {
		title.SetText(existing.Title)
		statusSelect.SetSelected(existing.Status)
		severitySelect.SetSelected(existing.Severity)
		dateNoticed.SetText(existing.DateNoticed.Format("2006-01-02"))
		location.SetText(existing.Location)
		appSelect.SetSelectedIndex(0)
		if existing.ApplianceID != nil {
			for i, id := range appIDs {
				if id == *existing.ApplianceID {
					appSelect.SetSelectedIndex(i)
					break
				}
			}
		}
		vendorSelect.SetSelectedIndex(0)
		if existing.VendorID != nil {
			for i, id := range vendorIDs {
				if id == *existing.VendorID {
					vendorSelect.SetSelectedIndex(i)
					break
				}
			}
		}
		if existing.CostCents != nil {
			cost.SetText(fmt.Sprintf("%.2f", float64(*existing.CostCents)/100))
		}
		description.SetText(existing.Description)
		notes.SetText(existing.Notes)
	} else {
		statusSelect.SetSelected(data.IncidentStatusOpen)
		severitySelect.SetSelected(data.IncidentSeveritySoon)
		dateNoticed.SetText(time.Now().Format("2006-01-02"))
		appSelect.SetSelectedIndex(0)
		vendorSelect.SetSelectedIndex(0)
	}

	form := widget.NewForm(
		widget.NewFormItem("Title *", title),
		widget.NewFormItem("Status", statusSelect),
		widget.NewFormItem("Severity", severitySelect),
		widget.NewFormItem("Date Noticed", dateNoticed),
		widget.NewFormItem("Location", location),
		widget.NewFormItem("Appliance", appSelect),
		widget.NewFormItem("Vendor", vendorSelect),
		widget.NewFormItem("Cost ($)", cost),
		widget.NewFormItem("Description", description),
		widget.NewFormItem("Notes", notes),
	)

	label := "Add Incident"
	if existing != nil {
		label = "Edit Incident"
	}

	d := dialog.NewCustomConfirm(label, "Save", "Cancel", form,
		func(confirmed bool) {
			if !confirmed {
				return
			}
			if title.Text == "" {
				dialog.ShowInformation("Validation", "Title is required.", w)
				return
			}

			dateN := time.Now()
			if t := parseDateOr(dateNoticed.Text); t != nil {
				dateN = *t
			}

			appIdx := appSelect.SelectedIndex()
			var appID *string
			if appIdx > 0 {
				id := appIDs[appIdx]
				appID = &id
			}

			vendorIdx := vendorSelect.SelectedIndex()
			var vendorID *string
			if vendorIdx > 0 {
				id := vendorIDs[vendorIdx]
				vendorID = &id
			}

			var costCents *int64
			if cost.Text != "" {
				v := parseFloatOr(cost.Text, -1)
				if v >= 0 {
					c := int64(v * 100)
					costCents = &c
				}
			}

			if existing != nil {
				existing.Title = title.Text
				existing.Status = statusSelect.Selected
				existing.Severity = severitySelect.Selected
				existing.DateNoticed = dateN
				existing.Location = location.Text
				existing.ApplianceID = appID
				existing.VendorID = vendorID
				existing.CostCents = costCents
				existing.Description = description.Text
				existing.Notes = notes.Text
				if err := store.UpdateIncident(*existing); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				now := time.Now()
				inc := data.Incident{
					ID:          uid.New(),
					Title:       title.Text,
					Status:      statusSelect.Selected,
					Severity:    severitySelect.Selected,
					DateNoticed: dateN,
					Location:    location.Text,
					ApplianceID: appID,
					VendorID:    vendorID,
					CostCents:   costCents,
					Description: description.Text,
					Notes:       notes.Text,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				if err := store.CreateIncident(&inc); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}
			if onSave != nil {
				onSave()
			}
		}, w)
	d.Resize(fyne.NewSize(520, 560))
	d.Show()
}

func incidentStatuses() []string {
	return []string{
		data.IncidentStatusOpen,
		data.IncidentStatusInProgress,
		data.IncidentStatusResolved,
	}
}

func incidentSeverities() []string {
	return []string{
		data.IncidentSeverityUrgent,
		data.IncidentSeveritySoon,
		data.IncidentSeverityWhenever,
	}
}
