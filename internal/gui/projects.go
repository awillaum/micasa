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

func newProjectsTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	var projects []data.Project
	var selected int = -1

	table := widget.NewTable(
		func() (int, int) { return len(projects) + 1, 5 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, cell fyne.CanvasObject) {
			lbl := cell.(*widget.Label)
			if id.Row == 0 {
				switch id.Col {
				case 0:
					lbl.SetText("Title")
				case 1:
					lbl.SetText("Type")
				case 2:
					lbl.SetText("Status")
				case 3:
					lbl.SetText("Budget")
				case 4:
					lbl.SetText("Updated")
				}
				lbl.TextStyle = fyne.TextStyle{Bold: true}
				return
			}
			i := id.Row - 1
			if i >= len(projects) {
				lbl.SetText("")
				return
			}
			p := projects[i]
			switch id.Col {
			case 0:
				lbl.SetText(p.Title)
			case 1:
				lbl.SetText(p.ProjectType.Name)
			case 2:
				lbl.SetText(p.Status)
			case 3:
				lbl.SetText(centsStr(p.BudgetCents))
			case 4:
				lbl.SetText(p.UpdatedAt.Format("2006-01-02"))
			}
		},
	)
	table.SetColumnWidth(0, 250)
	table.SetColumnWidth(1, 120)
	table.SetColumnWidth(2, 100)
	table.SetColumnWidth(3, 100)
	table.SetColumnWidth(4, 100)

	table.OnSelected = func(id widget.TableCellID) {
		if id.Row > 0 {
			selected = id.Row - 1
		}
	}

	reload := func() {
		var err error
		projects, err = store.ListProjects(false)
		if err != nil {
			dialog.ShowError(err, w)
			return
		}
		selected = -1
		table.Refresh()
	}
	reload()

	addBtn := widget.NewButton("Add", func() {
		showProjectDialog(store, w, nil, reload)
	})

	editBtn := widget.NewButton("Edit", func() {
		if selected < 0 || selected >= len(projects) {
			dialog.ShowInformation("No selection", "Please select a project to edit.", w)
			return
		}
		p := projects[selected]
		showProjectDialog(store, w, &p, reload)
	})

	deleteBtn := widget.NewButton("Delete", func() {
		if selected < 0 || selected >= len(projects) {
			dialog.ShowInformation("No selection", "Please select a project to delete.", w)
			return
		}
		p := projects[selected]
		dialog.ShowConfirm("Delete Project",
			fmt.Sprintf("Delete project %q?", p.Title),
			func(ok bool) {
				if !ok {
					return
				}
				if err := store.DeleteProject(p.ID); err != nil {
					dialog.ShowError(err, w)
					return
				}
				reload()
			}, w)
	})

	toolbar := container.NewHBox(addBtn, editBtn, deleteBtn)
	return container.NewBorder(toolbar, nil, nil, nil, table)
}

func showProjectDialog(store *data.Store, w fyne.Window, existing *data.Project, onSave func()) {
	types, err := store.ProjectTypes()
	if err != nil {
		dialog.ShowError(err, w)
		return
	}

	typeNames := make([]string, len(types))
	typeIDs := make([]string, len(types))
	for i, t := range types {
		typeNames[i] = t.Name
		typeIDs[i] = t.ID
	}

	title := widget.NewEntry()
	typeSelect := widget.NewSelect(typeNames, nil)
	statusSelect := widget.NewSelect(projectStatuses(), nil)
	budget := widget.NewEntry()
	budget.SetPlaceHolder("e.g. 1500.00")
	description := widget.NewMultiLineEntry()
	description.SetMinRowsVisible(3)

	if existing != nil {
		title.SetText(existing.Title)
		for i, id := range typeIDs {
			if id == existing.ProjectTypeID {
				typeSelect.SetSelectedIndex(i)
				break
			}
		}
		statusSelect.SetSelected(existing.Status)
		if existing.BudgetCents != nil {
			budget.SetText(fmt.Sprintf("%.2f", float64(*existing.BudgetCents)/100))
		}
		description.SetText(existing.Description)
	} else if len(typeNames) > 0 {
		typeSelect.SetSelectedIndex(0)
		statusSelect.SetSelected(data.ProjectStatusPlanned)
	}

	form := widget.NewForm(
		widget.NewFormItem("Title *", title),
		widget.NewFormItem("Type *", typeSelect),
		widget.NewFormItem("Status", statusSelect),
		widget.NewFormItem("Budget ($)", budget),
		widget.NewFormItem("Description", description),
	)

	label := "Add Project"
	if existing != nil {
		label = "Edit Project"
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
			typeIdx := typeSelect.SelectedIndex()
			if typeIdx < 0 {
				dialog.ShowInformation("Validation", "Project type is required.", w)
				return
			}

			var budgetCents *int64
			if budget.Text != "" {
				v := parseFloatOr(budget.Text, -1)
				if v < 0 {
					dialog.ShowInformation("Validation", "Budget must be a positive number.", w)
					return
				}
				c := int64(v * 100)
				budgetCents = &c
			}

			if existing != nil {
				existing.Title = title.Text
				existing.ProjectTypeID = typeIDs[typeIdx]
				existing.Status = statusSelect.Selected
				existing.BudgetCents = budgetCents
				existing.Description = description.Text
				if err := store.UpdateProject(*existing); err != nil {
					dialog.ShowError(err, w)
					return
				}
			} else {
				now := time.Now()
				p := data.Project{
					ID:            uid.New(),
					Title:         title.Text,
					ProjectTypeID: typeIDs[typeIdx],
					Status:        statusSelect.Selected,
					BudgetCents:   budgetCents,
					Description:   description.Text,
					CreatedAt:     now,
					UpdatedAt:     now,
				}
				if err := store.CreateProject(&p); err != nil {
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

func projectStatuses() []string {
	return []string{
		data.ProjectStatusIdeating,
		data.ProjectStatusPlanned,
		data.ProjectStatusQuoted,
		data.ProjectStatusInProgress,
		data.ProjectStatusDelayed,
		data.ProjectStatusCompleted,
		data.ProjectStatusAbandoned,
	}
}
