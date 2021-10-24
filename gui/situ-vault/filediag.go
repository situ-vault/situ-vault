package main

import (
	"os"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
)

func wdListUri() (fyne.ListableURI, error) {
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	listUri, err := storage.ListerForURI(storage.NewFileURI(wd))
	return listUri, err
}
