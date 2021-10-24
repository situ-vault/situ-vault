package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
)

var driver fyne.Driver

func main() {
	application := app.New()
	driver = application.Driver()
	exp := newExperience()
	exp.loadUi(application)
	application.Run()
}
