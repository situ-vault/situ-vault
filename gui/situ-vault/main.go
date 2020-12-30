package main

import (
	"fyne.io/fyne"
	"fyne.io/fyne/app"
)

var driver fyne.Driver

func main() {
	application := app.New()
	driver = application.Driver()
	exp := newExperience()
	exp.loadUi(application)
	application.Run()
}
