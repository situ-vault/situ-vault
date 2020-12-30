package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
)

type ToolbarLabeledAction struct {
	Icon        fyne.Resource
	Label       string
	OnActivated func()
	resButton   *widget.Button
}

func NewToolbarLabeledAction(icon fyne.Resource, label string, onActivated func()) ToolbarLabeledItem {
	button := widget.NewButtonWithIcon(label, icon, onActivated)
	button.Importance = widget.LowImportance
	action := ToolbarLabeledAction{icon, label, onActivated, button}
	return &action
}

func (t ToolbarLabeledAction) ToolbarObject() fyne.CanvasObject {
	return t.resButton
}

type ToolbarLabeledItem interface {
	ToolbarObject() fyne.CanvasObject
	Button() *widget.Button
}

func (t *ToolbarLabeledAction) Button() *widget.Button {
	return t.resButton
}
