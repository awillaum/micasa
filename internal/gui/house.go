// Copyright 2026 Phillip Cloud
// Licensed under the Apache License, Version 2.0
//go:build gui


package gui

import (
	"fmt"
	"strconv"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/micasa-dev/micasa/internal/data"
)

func newHouseTab(store *data.Store, w fyne.Window) fyne.CanvasObject {
	profile, err := store.HouseProfile()
	if err != nil {
		return widget.NewLabel("No house profile found. Use the Edit button to create one.")
	}

	form := widget.NewForm(houseFormItems(profile)...)

	editBtn := widget.NewButton("Edit", func() {
		showHouseDialog(store, w, func() {
			// Refresh after save.
			updated, err2 := store.HouseProfile()
			if err2 == nil {
				form.Items = houseFormItems(updated)
				form.Refresh()
			}
		})
	})

	return container.NewBorder(
		container.NewHBox(editBtn),
		nil, nil, nil,
		container.NewScroll(form),
	)
}

func houseFormItems(p data.HouseProfile) []*widget.FormItem {
	return []*widget.FormItem{
		widget.NewFormItem("Nickname", widget.NewLabel(p.Nickname)),
		widget.NewFormItem("Address", widget.NewLabel(formatAddress(p))),
		widget.NewFormItem("Year Built", widget.NewLabel(intStr(p.YearBuilt))),
		widget.NewFormItem("Sq Ft", widget.NewLabel(intStr(p.SquareFeet))),
		widget.NewFormItem("Lot Sq Ft", widget.NewLabel(intStr(p.LotSquareFeet))),
		widget.NewFormItem("Bedrooms", widget.NewLabel(intStr(p.Bedrooms))),
		widget.NewFormItem("Bathrooms", widget.NewLabel(fmt.Sprintf("%.1f", p.Bathrooms))),
		widget.NewFormItem("Foundation", widget.NewLabel(p.FoundationType)),
		widget.NewFormItem("Wiring", widget.NewLabel(p.WiringType)),
		widget.NewFormItem("Roof", widget.NewLabel(p.RoofType)),
		widget.NewFormItem("Exterior", widget.NewLabel(p.ExteriorType)),
		widget.NewFormItem("Heating", widget.NewLabel(p.HeatingType)),
		widget.NewFormItem("Cooling", widget.NewLabel(p.CoolingType)),
		widget.NewFormItem("Water", widget.NewLabel(p.WaterSource)),
		widget.NewFormItem("Sewer", widget.NewLabel(p.SewerType)),
		widget.NewFormItem("Parking", widget.NewLabel(p.ParkingType)),
		widget.NewFormItem("Basement", widget.NewLabel(p.BasementType)),
		widget.NewFormItem("Insurance Carrier", widget.NewLabel(p.InsuranceCarrier)),
		widget.NewFormItem("Insurance Policy", widget.NewLabel(p.InsurancePolicy)),
		widget.NewFormItem("HOA Name", widget.NewLabel(p.HOAName)),
	}
}

func formatAddress(p data.HouseProfile) string {
	addr := p.AddressLine1
	if p.AddressLine2 != "" {
		addr += ", " + p.AddressLine2
	}
	if p.City != "" {
		addr += ", " + p.City
	}
	if p.State != "" {
		addr += ", " + p.State
	}
	if p.PostalCode != "" {
		addr += " " + p.PostalCode
	}
	return addr
}

func intStr(v int) string {
	if v == 0 {
		return ""
	}
	return strconv.Itoa(v)
}

func showHouseDialog(store *data.Store, w fyne.Window, onSave func()) {
	profile, _ := store.HouseProfile()

	nickname := widget.NewEntry()
	nickname.SetText(profile.Nickname)
	addr1 := widget.NewEntry()
	addr1.SetText(profile.AddressLine1)
	addr2 := widget.NewEntry()
	addr2.SetText(profile.AddressLine2)
	city := widget.NewEntry()
	city.SetText(profile.City)
	state := widget.NewEntry()
	state.SetText(profile.State)
	postal := widget.NewEntry()
	postal.SetText(profile.PostalCode)
	yearBuilt := widget.NewEntry()
	yearBuilt.SetText(intStr(profile.YearBuilt))
	sqft := widget.NewEntry()
	sqft.SetText(intStr(profile.SquareFeet))
	lotSqft := widget.NewEntry()
	lotSqft.SetText(intStr(profile.LotSquareFeet))
	bedrooms := widget.NewEntry()
	bedrooms.SetText(intStr(profile.Bedrooms))
	bathrooms := widget.NewEntry()
	bathrooms.SetText(fmt.Sprintf("%.1f", profile.Bathrooms))
	foundation := widget.NewEntry()
	foundation.SetText(profile.FoundationType)
	wiring := widget.NewEntry()
	wiring.SetText(profile.WiringType)
	roof := widget.NewEntry()
	roof.SetText(profile.RoofType)
	exterior := widget.NewEntry()
	exterior.SetText(profile.ExteriorType)
	heating := widget.NewEntry()
	heating.SetText(profile.HeatingType)
	cooling := widget.NewEntry()
	cooling.SetText(profile.CoolingType)
	water := widget.NewEntry()
	water.SetText(profile.WaterSource)
	sewer := widget.NewEntry()
	sewer.SetText(profile.SewerType)
	parking := widget.NewEntry()
	parking.SetText(profile.ParkingType)
	basement := widget.NewEntry()
	basement.SetText(profile.BasementType)
	insuranceCarrier := widget.NewEntry()
	insuranceCarrier.SetText(profile.InsuranceCarrier)
	insurancePolicy := widget.NewEntry()
	insurancePolicy.SetText(profile.InsurancePolicy)
	hoaName := widget.NewEntry()
	hoaName.SetText(profile.HOAName)

	form := widget.NewForm(
		widget.NewFormItem("Nickname", nickname),
		widget.NewFormItem("Address Line 1", addr1),
		widget.NewFormItem("Address Line 2", addr2),
		widget.NewFormItem("City", city),
		widget.NewFormItem("State", state),
		widget.NewFormItem("Postal Code", postal),
		widget.NewFormItem("Year Built", yearBuilt),
		widget.NewFormItem("Sq Ft", sqft),
		widget.NewFormItem("Lot Sq Ft", lotSqft),
		widget.NewFormItem("Bedrooms", bedrooms),
		widget.NewFormItem("Bathrooms", bathrooms),
		widget.NewFormItem("Foundation", foundation),
		widget.NewFormItem("Wiring", wiring),
		widget.NewFormItem("Roof", roof),
		widget.NewFormItem("Exterior", exterior),
		widget.NewFormItem("Heating", heating),
		widget.NewFormItem("Cooling", cooling),
		widget.NewFormItem("Water Source", water),
		widget.NewFormItem("Sewer", sewer),
		widget.NewFormItem("Parking", parking),
		widget.NewFormItem("Basement", basement),
		widget.NewFormItem("Insurance Carrier", insuranceCarrier),
		widget.NewFormItem("Insurance Policy", insurancePolicy),
		widget.NewFormItem("HOA Name", hoaName),
	)

	d := dialog.NewCustomConfirm("Edit House Profile", "Save", "Cancel",
		container.NewScroll(form),
		func(confirmed bool) {
			if !confirmed {
				return
			}
			profile.Nickname = nickname.Text
			profile.AddressLine1 = addr1.Text
			profile.AddressLine2 = addr2.Text
			profile.City = city.Text
			profile.State = state.Text
			profile.PostalCode = postal.Text
			profile.YearBuilt = parseIntOr(yearBuilt.Text, 0)
			profile.SquareFeet = parseIntOr(sqft.Text, 0)
			profile.LotSquareFeet = parseIntOr(lotSqft.Text, 0)
			profile.Bedrooms = parseIntOr(bedrooms.Text, 0)
			profile.Bathrooms = parseFloatOr(bathrooms.Text, 0)
			profile.FoundationType = foundation.Text
			profile.WiringType = wiring.Text
			profile.RoofType = roof.Text
			profile.ExteriorType = exterior.Text
			profile.HeatingType = heating.Text
			profile.CoolingType = cooling.Text
			profile.WaterSource = water.Text
			profile.SewerType = sewer.Text
			profile.ParkingType = parking.Text
			profile.BasementType = basement.Text
			profile.InsuranceCarrier = insuranceCarrier.Text
			profile.InsurancePolicy = insurancePolicy.Text
			profile.HOAName = hoaName.Text

			var saveErr error
			if profile.ID == "" {
				saveErr = store.CreateHouseProfile(profile)
			} else {
				saveErr = store.UpdateHouseProfile(profile)
			}
			if saveErr != nil {
				dialog.ShowError(saveErr, w)
				return
			}
			if onSave != nil {
				onSave()
			}
		},
		w,
	)
	d.Resize(fyne.NewSize(600, 500))
	d.Show()
}
