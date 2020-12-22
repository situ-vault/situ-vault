package main

import "fyne.io/fyne/app"

func main() {
	application := app.New()
	exp := newExperience()
	exp.loadUi(application)
	application.Run()
}
