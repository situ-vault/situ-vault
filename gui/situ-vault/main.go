package main

import (
	"fyne.io/fyne/v2/app"
)

func main() {
	application := app.New()
	application.Settings().SetTheme(&fyneTheme{})
	_ = application.Driver()
	exp := newExperience()
	exp.loadUi(application)
	application.Run()
}
