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

func newAppliancesTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var appliances []data.Appliance
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(appliances) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Name")
				case 1:
					lbl.SetText("Brand")
				case 2:
					lbl.SetText("Model")
				case 3:
					lbl.SetText("Location")
				case 4:
					lbl.SetText("Warranty")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(appliances) {
				lbl.SetText("")
				return
			}
			a := appliances[i]
			switch id.Col {
			case 0:
				lbl.SetText(a.Name)
			case 1:
				lbl.SetText(a.Brand)
			case 2:
				lbl.SetText(a.ModelNumber)
			case 3:
				lbl.SetText(a.Location)
			case 4:
				lbl.SetText(timeStr(a.WarrantyExpiry))
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 130)
	table.SetColumnWidth(2, 130)
	table.SetColumnWidth(3, 150)
	table.SetColumnWidth(4, 110)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		appliances, err = store.ListAppliances(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showApplianceDialog(store, w, nil, reload)
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(appliances) {
			dialog.ShowInformation("No selection", "Please select an appliance to edit.", w)
			return
		}
		a := appliances[selected]
		showApplianceDialog(store, w, &a, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(appliances) {
			dialog.ShowInformation("No selection", "Please select an appliance to delete.", w)
			return
		}
		a := appliances[selected]
		dialog.ShowConfirm("Delete Appliance",
			fmt.Sprintf("Delete appliance %q?", a.Name),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteAppliance(a.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showApplianceDialog(store *data.Store, w fyne.Window, existing *data.Appliance, onSave func()) {
	name := widget.NewEntry()
	brand := widget.NewEntry()
	modelNum := widget.NewEntry()
	serial := widget.NewEntry()
	location := widget.NewEntry()
	purchase := widget.NewEntry()
	purchase.SetPlaceHolder("YYYY-MM-DD")
	warranty := widget.NewEntry()
	warranty.SetPlaceHolder("YYYY-MM-DD")
	cost := widget.NewEntry()
	cost.SetPlaceHolder("e.g. 500.00")
	notes := widget.NewMultiLineEntry()
	notes.SetMinRowsVisible(3)

	if existing != nil {
		name.SetText(existing.Name)
		brand.SetText(existing.Brand)
		modelNum.SetText(existing.ModelNumber)
		serial.SetText(existing.SerialNumber)
		location.SetText(existing.Location)
		purchase.SetText(timeStr(existing.PurchaseDate))
		warranty.SetText(timeStr(existing.WarrantyExpiry))
		if existing.CostCents != nil {
			cost.SetText(fmt.Sprintf("%.2f", float64(*existing.CostCents)/100))
		}
		notes.SetText(existing.Notes)
	}

	form := widget.NewForm(
		widget.NewFormItem("Name *", name),
		widget.NewFormItem("Brand", brand),
		widget.NewFormItem("Model Number", modelNum),
		widget.NewFormItem("Serial Number", serial),
		widget.NewFormItem("Location", location),
		widget.NewFormItem("Purchase Date", purchase),
		widget.NewFormItem("Warranty Expiry", warranty),
		widget.NewFormItem("Cost ($)", cost),
		widget.NewFormItem("Notes", notes),
	)

	label := "Add Appliance"
	if existing != nil {
		label = "Edit Appliance"
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

			purchaseDate := parseDateOr(purchase.Text)
			warrantyDate := parseDateOr(warranty.Text)
			var costCents *int64
			if cost.Text != "" {
				v := parseFloatOr(cost.Text, -1)
				if v >= 0 {
					c := int64(v * 100)
					costCents = &c
				}
			}

			if existing != nil {
				existing.Name = name.Text
				existing.Brand = brand.Text
				existing.ModelNumber = modelNum.Text
				existing.SerialNumber = serial.Text
				existing.Location = location.Text
				existing.PurchaseDate = purchaseDate
				existing.WarrantyExpiry = warrantyDate
				existing.CostCents = costCents
				existing.Notes = notes.Text
				if err := store.UpdateAppliance(*existing); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				now := time.Now()
				a := data.Appliance{
					ID:             uid.New(),
					Name:           name.Text,
					Brand:          brand.Text,
					ModelNumber:    modelNum.Text,
					SerialNumber:   serial.Text,
					Location:       location.Text,
					PurchaseDate:   purchaseDate,
					WarrantyExpiry: warrantyDate,
					CostCents:      costCents,
					Notes:          notes.Text,
					CreatedAt:      now,
					UpdatedAt:      now,
				}
				if err := store.CreateAppliance(&a); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}
			if onSave != nil {
				onSave()
			}
		}, w)
	d.Resize(fyne.NewSize(500, 450))
	d.Show()
}
