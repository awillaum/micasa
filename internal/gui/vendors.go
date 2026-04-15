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

func newVendorsTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var vendors []data.Vendor
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(vendors) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Name")
				case 1:
					lbl.SetText("Contact")
				case 2:
					lbl.SetText("Phone")
				case 3:
					lbl.SetText("Email")
				case 4:
					lbl.SetText("Website")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(vendors) {
				lbl.SetText("")
				return
			}
			v := vendors[i]
			switch id.Col {
			case 0:
				lbl.SetText(v.Name)
			case 1:
				lbl.SetText(v.ContactName)
			case 2:
				lbl.SetText(v.Phone)
			case 3:
				lbl.SetText(v.Email)
			case 4:
				lbl.SetText(v.Website)
			}
		},
	)
	table.SetColumnWidth(0, 200)
	table.SetColumnWidth(1, 150)
	table.SetColumnWidth(2, 130)
	table.SetColumnWidth(3, 200)
	table.SetColumnWidth(4, 200)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		vendors, err = store.ListVendors(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showVendorDialog(store, w, nil, reload)
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(vendors) {
			dialog.ShowInformation("No selection", "Please select a vendor to edit.", w)
			return
		}
		v := vendors[selected]
		showVendorDialog(store, w, &v, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(vendors) {
			dialog.ShowInformation("No selection", "Please select a vendor to delete.", w)
			return
		}
		v := vendors[selected]
		dialog.ShowConfirm("Delete Vendor",
			fmt.Sprintf("Delete vendor %q?", v.Name),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteVendor(v.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showVendorDialog(store *data.Store, w fyne.Window, existing *data.Vendor, onSave func()) {
	name := widget.NewEntry()
	contact := widget.NewEntry()
	phone := widget.NewEntry()
	email := widget.NewEntry()
	website := widget.NewEntry()
	notes := widget.NewMultiLineEntry()
	notes.SetMinRowsVisible(3)

	if existing != nil {
		name.SetText(existing.Name)
		contact.SetText(existing.ContactName)
		phone.SetText(existing.Phone)
		email.SetText(existing.Email)
		website.SetText(existing.Website)
		notes.SetText(existing.Notes)
	}

	form := widget.NewForm(
		widget.NewFormItem("Name *", name),
		widget.NewFormItem("Contact Name", contact),
		widget.NewFormItem("Phone", phone),
		widget.NewFormItem("Email", email),
		widget.NewFormItem("Website", website),
		widget.NewFormItem("Notes", notes),
	)

	label := "Add Vendor"
	if existing != nil {
		label = "Edit Vendor"
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

			if existing != nil {
				existing.Name = name.Text
				existing.ContactName = contact.Text
				existing.Phone = phone.Text
				existing.Email = email.Text
				existing.Website = website.Text
				existing.Notes = notes.Text
				if err := store.UpdateVendor(*existing); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				now := time.Now()
				v := data.Vendor{
					ID:          uid.New(),
					Name:        name.Text,
					ContactName: contact.Text,
					Phone:       phone.Text,
					Email:       email.Text,
					Website:     website.Text,
					Notes:       notes.Text,
					CreatedAt:   now,
					UpdatedAt:   now,
				}
				if err := store.CreateVendor(&v); err != nil {
					dialog.ShowError(err, w)
					return
				}
			}
			if onSave != nil {
				onSave()
			}
		}, w)
	d.Resize(fyne.NewSize(500, 400))
	d.Show()
}
