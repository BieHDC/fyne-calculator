// Package main launches the calculator app
//
//go:generate fyne bundle -o data.go Icon.png
package main

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	
	"github.com/BieHDC/fyne-calculator/calc"
)

func main() {
	app := app.New()
	app.SetIcon(resourceIconPng)

	window := app.NewWindow("Calc")
	c := calc.NewCalculator()
	c.ConnectKeyboard(window)
	window.SetContent(c.Content())

	window.Resize(fyne.NewSize(200, 300))
	window.ShowAndRun()
}
