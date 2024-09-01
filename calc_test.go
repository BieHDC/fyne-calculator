package main

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"

	"github.com/stretchr/testify/assert"
)

func TestAdd(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "2", calc.output.Text)
}

func TestSubtract(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["-"])
	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "1", calc.output.Text)
}

func TestDivide(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["3"])
	test.Tap(calc.buttons["/"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "1.5", calc.output.Text)
}

func TestMultiply(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["5"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "10", calc.output.Text)
}

func TestParenthesis(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["*"])
	test.Tap(calc.buttons["("])
	test.Tap(calc.buttons["3"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["4"])
	test.Tap(calc.buttons[")"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "14", calc.output.Text)
}

func TestDot(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["7"])
	test.Tap(calc.buttons["."])
	test.Tap(calc.buttons["8"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "10", calc.output.Text)
}

func TestClear(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["1"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["C"])

	assert.Equal(t, "", calc.output.Text)
}

func TestContinueAfterResult(t *testing.T) {
	calc := newCalculator()
	calc.ConnectKeyboard(test.NewApp().NewWindow(""))

	test.Tap(calc.buttons["6"])
	test.Tap(calc.buttons["+"])
	test.Tap(calc.buttons["4"])
	test.Tap(calc.buttons["="])
	test.Tap(calc.buttons["-"])
	test.Tap(calc.buttons["2"])
	test.Tap(calc.buttons["="])

	assert.Equal(t, "8", calc.output.Text)
}

func TestKeyboard(t *testing.T) {
	calc := newCalculator()
	window := test.NewApp().NewWindow("")
	calc.ConnectKeyboard(window)

	test.TypeOnCanvas(window.Canvas(), "1+1")
	assert.Equal(t, "1+1", calc.output.Text)

	test.TypeOnCanvas(window.Canvas(), "=")
	assert.Equal(t, "2", calc.output.Text)

	test.TypeOnCanvas(window.Canvas(), "c")
	assert.Equal(t, "", calc.output.Text)
}

func TestKeyboard_Buttons(t *testing.T) {
	calc := newCalculator()
	window := test.NewApp().NewWindow("")
	calc.ConnectKeyboard(window)

	test.TypeOnCanvas(window.Canvas(), "1+1")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyReturn})
	assert.Equal(t, "2", calc.output.Text)

	test.TypeOnCanvas(window.Canvas(), "c")

	test.TypeOnCanvas(window.Canvas(), "1+1")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	assert.Equal(t, "2", calc.output.Text)
}

func TestKeyboard_Backspace(t *testing.T) {
	calc := newCalculator()
	window := test.NewApp().NewWindow("")
	calc.ConnectKeyboard(window)

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "", calc.output.Text)

	test.TypeOnCanvas(window.Canvas(), "1/2")
	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "1/", calc.output.Text)

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyEnter})
	assert.Equal(t, "Unexpected end of expression", calc.errline.Content.(*widget.Label).Text)

	calc.onTypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	assert.Equal(t, "Unexpected end of expression", calc.errline.Content.(*widget.Label).Text)
}

func TestError(t *testing.T) {
	calc := newCalculator()
	window := test.NewApp().NewWindow("")
	calc.ConnectKeyboard(window)

	test.TypeOnCanvas(window.Canvas(), "1//1=")
	assert.Equal(t, "Invalid token: '//'", calc.errline.Content.(*widget.Label).Text)

	test.TypeOnCanvas(window.Canvas(), "c")

	test.TypeOnCanvas(window.Canvas(), "()9=")
	assert.Equal(t, "Input cant be float64ed: <nil>", calc.errline.Content.(*widget.Label).Text)

	test.TypeOnCanvas(window.Canvas(), "=")
	assert.Equal(t, "Input cant be float64ed: <nil>", calc.errline.Content.(*widget.Label).Text)

	test.TypeOnCanvas(window.Canvas(), "55=")
	assert.Equal(t, "Input cant be float64ed: <nil>", calc.errline.Content.(*widget.Label).Text)
}

func TestShortcuts(t *testing.T) {
	app := test.NewApp()
	calc := newCalculator()
	window := app.NewWindow("")
	calc.ConnectKeyboard(window)
	clipboard := window.Clipboard()

	test.TypeOnCanvas(window.Canvas(), "720 + 80")
	calc.onCopyShortcut(&fyne.ShortcutCopy{Clipboard: clipboard})
	assert.Equal(t, clipboard.Content(), calc.output.Text)

	test.TypeOnCanvas(window.Canvas(), "+")
	clipboard.SetContent("50")
	calc.onPasteShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	test.TypeOnCanvas(window.Canvas(), "=")
	assert.Equal(t, "850", calc.output.Text)

	clipboard.SetContent("not a valid number")
	calc.onPasteShortcut(&fyne.ShortcutPaste{Clipboard: clipboard})
	assert.Equal(t, "850", calc.output.Text)
}
