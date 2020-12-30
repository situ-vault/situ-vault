package main

import (
	"os"

	"fyne.io/fyne"
	"fyne.io/fyne/storage"
)

func wdListUri() (fyne.ListableURI, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	listUri, err := driver.ListerForURI(storage.NewFileURI(wd))
	return listUri, err
}
