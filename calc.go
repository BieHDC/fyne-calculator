//go:generate fyne bundle -o data.go Icon.png

package main

import (
	"errors"
	"fmt"
	"strconv"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/widget"

	"github.com/Knetic/govaluate"
)

type calc struct {
	output  *widget.Entry
	errline *container.Scroll

	buttons map[string]*widget.Button

	content fyne.CanvasObject
}

func (c *calc) typeKeys(s string) {
	for _, r := range s {
		c.output.TypedRune(r)
	}
}

func (c *calc) character(char rune) {
	c.output.TypedRune(char)
}

func (c *calc) digit(d int) {
	c.character(rune(d) + '0')
}

func (c *calc) addButton(text string, action func()) *widget.Button {
	button := widget.NewButton(text, action)
	c.buttons[text] = button
	return button
}

func (c *calc) digitButton(number int) *widget.Button {
	str := strconv.Itoa(number)
	return c.addButton(str, func() {
		c.digit(number)
	})
}

func (c *calc) charButton(char rune) *widget.Button {
	return c.addButton(string(char), func() {
		c.character(char)
	})
}

func validRune(r rune) bool {
	switch r {
	// numbers
	case '1', '2', '3', '4', '5', '6', '7', '8', '9', '0':
		return true
	// operators
	case '(', ')', '/', '*', '-', '+', ',', '.':
		return true
	// more things the arith lib supports and i will allow
	case '&', '|', '^', '%', '>', '<', '!', '~', '?', ':', '=':
		return true
	// nopers
	default:
		return false
	}
}

func checkInput(s string) (string, error) {
	for i, r := range s {
		if validRune(r) {
			continue
		}

		// we have something not valid
		// get a run of it and return it as error
		var j int
		var rr rune
		for j, rr = range s[i:] {
			if validRune(rr) {
				break
			}
		}
		// give the user a hint where
		return "", errors.New(s[i : i+j+1])
	}

	return s, nil
}

func (c *calc) evaluate() {
	// sanitise input because for reasons the thing panics otherwise
	// one thing the library supports and i dont is arrays since we
	// replace ',' with '.' due to input.
	sanitised, err := checkInput(strings.ReplaceAll(c.output.Text, ",", "."))
	if err != nil {
		c.setErrline("Invalid input at: "+err.Error(), true)
		return
	}

	expression, err := govaluate.NewEvaluableExpression(sanitised)
	if err != nil {
		c.setErrline(err.Error(), true)
		return
	}

	result, err := expression.Evaluate(nil)
	if err != nil {
		c.setErrline(err.Error(), true)
		return
	}

	value, ok := result.(float64)
	if !ok {
		c.setErrline(fmt.Sprintf("Input cant be float64ed: %v", result), true)
		return
	}

	c.setTextWithUndoPreserve(strconv.FormatFloat(value, 'f', -1, 64))
	c.setErrline("", false) // no error
}

func (c *calc) setErrline(s string, show bool) {
	c.errline.Content.(*widget.Label).SetText(s)
	if show {
		c.errline.ScrollToTop()
		c.errline.Show()
	} else {
		c.errline.Hide()
	}
}

// all it does is select the current line and then
// type in the new text to preserve the undo/redo
func (c *calc) setTextWithUndoPreserve(s string) {
	c.output.DoubleTapped(&fyne.PointEvent{}) // double tap as prep
	c.output.MouseDown(&desktop.MouseEvent{}) // and now troll it into a triple tap to select everything
	if len(s) > 0 {
		// and type out the result
		c.typeKeys(s)
	} else {
		// or clear the output
		c.output.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	}
}

func (c *calc) onTypedRune(r rune) {
	if r == 'c' {
		r = 'C' // The button is using a capital C.
	}

	if button, ok := c.buttons[string(r)]; ok {
		button.OnTapped()
	}
}

func (c *calc) onTypedKey(ev *fyne.KeyEvent) {
	if ev.Name == fyne.KeyReturn || ev.Name == fyne.KeyEnter {
		c.evaluate()
	} else if ev.Name == fyne.KeyBackspace {
		c.output.TypedKey(&fyne.KeyEvent{Name: fyne.KeyBackspace})
	}
}

func (c *calc) onPasteShortcut(shortcut fyne.Shortcut) {
	content := shortcut.(*fyne.ShortcutPaste).Clipboard.Content()
	if _, err := strconv.ParseFloat(content, 64); err != nil {
		return
	}

	c.typeKeys(content)
}

func (c *calc) onCopyShortcut(shortcut fyne.Shortcut) {
	shortcut.(*fyne.ShortcutCopy).Clipboard.SetContent(c.output.Text)
}

func (c *calc) ConnectKeyboard(window fyne.Window) {
	canvas := window.Canvas()
	canvas.SetOnTypedRune(c.onTypedRune)
	canvas.SetOnTypedKey(c.onTypedKey)
	canvas.AddShortcut(&fyne.ShortcutCopy{}, c.onCopyShortcut)
	canvas.AddShortcut(&fyne.ShortcutPaste{}, c.onPasteShortcut)
}

func newCalculator() *calc {
	var c calc

	c.buttons = make(map[string]*widget.Button, 19)

	c.output = widget.NewEntry()
	c.output.TextStyle = fyne.TextStyle{Monospace: true}
	c.output.OnSubmitted = func(_ string) { c.evaluate() }
	c.errline = container.NewHScroll(widget.NewLabel(""))
	c.errline.Hide()

	equals := c.addButton("=", c.evaluate)
	equals.Importance = widget.HighImportance

	c.content = container.NewGridWithColumns(1,
		c.output,
		c.errline,
		container.NewGridWithColumns(4,
			c.addButton("C", func() { c.setTextWithUndoPreserve("") }),
			c.charButton('('),
			c.charButton(')'),
			c.charButton('/')),
		container.NewGridWithColumns(4,
			c.digitButton(7),
			c.digitButton(8),
			c.digitButton(9),
			c.charButton('*')),
		container.NewGridWithColumns(4,
			c.digitButton(4),
			c.digitButton(5),
			c.digitButton(6),
			c.charButton('-')),
		container.NewGridWithColumns(4,
			c.digitButton(1),
			c.digitButton(2),
			c.digitButton(3),
			c.charButton('+')),
		container.NewGridWithColumns(2,
			container.NewGridWithColumns(2,
				c.digitButton(0),
				c.charButton('.')),
			equals),
	)

	// register the comma of other languages
	// fixme localise when required in future
	_ = c.addButton(",", func() { c.character('.') })

	return &c
}

func (c *calc) Content() fyne.CanvasObject {
	return c.content
}
