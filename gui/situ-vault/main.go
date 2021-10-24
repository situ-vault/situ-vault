package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var driver fyne.Driver

func main() {
	application := app.New()
	application.Settings().SetTheme(&fyneTheme{})
	driver = application.Driver()
	exp := newExperience()
	exp.loadUi(application)
	application.Run()
}
