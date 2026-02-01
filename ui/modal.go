package ui

import "github.com/rivo/tview"

func GetOKModal(app *tview.Application, switchToPrimitive tview.Primitive, text string) *tview.Modal {
	okModal := tview.NewModal().
		SetText(text).
		AddButtons([]string{"OK"}).
		SetDoneFunc(func(int, string) {
			app.SetRoot(switchToPrimitive, true)
		})
	return okModal
}
