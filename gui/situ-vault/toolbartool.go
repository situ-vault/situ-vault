package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/widget"
	"time"
)

// Action with text label in the Toolbar:

type ToolbarLabeledAction struct {
	Icon        fyne.Resource
	Label       string
	OnActivated func()
	resButton   *widget.Button
}

func NewToolbarLabeledAction(icon fyne.Resource, label string, onActivated func()) *ToolbarLabeledAction {
	button := widget.NewButtonWithIcon(label, icon, onActivated)
	button.Importance = widget.LowImportance
	action := &ToolbarLabeledAction{icon, label, onActivated, button}
	return action
}

func (t ToolbarLabeledAction) ToolbarObject() fyne.CanvasObject {
	return t.resButton
}

func (t ToolbarLabeledAction) Button() *widget.Button {
	return t.resButton
}

// Notification text shown in the Toolbar:

type ToolbarNotification struct {
	label *widget.Label
}

func (t ToolbarNotification) ToolbarObject() fyne.CanvasObject {
	return t.label
}

func (t ToolbarNotification) SetText(text string) {
	t.label.TextStyle.Italic = true
	t.label.SetText(text)
}

func NewToolbarNotification() *ToolbarNotification {
	tn := &ToolbarNotification{
		label: widget.NewLabel(""),
	}
	return tn
}

// ShowNotification show a notification and uses a goroutine to reset it after the sleep duration
func (t ToolbarNotification) ShowNotification(text string, sleep time.Duration, refresh func()) {
	t.SetText(text)
	refresh()
	resetLater := func() {
		time.Sleep(sleep)
		t.SetText("")
		refresh()
	}
	go resetLater()
}
