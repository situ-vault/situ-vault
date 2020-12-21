package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ToolbarLabeledAction struct {
	Icon        fyne.Resource
	Label 		string
	OnActivated func()
}

func NewToolbarLabeledAction(icon fyne.Resource, label string, onActivated func()) widget.ToolbarItem {
	return &ToolbarLabeledAction{icon,label, onActivated}
}

func (t ToolbarLabeledAction) ToolbarObject() fyne.CanvasObject {
	button := widget.NewButtonWithIcon(t.Label, t.Icon, t.OnActivated)
	button.Importance = widget.LowImportance
	return button
}
