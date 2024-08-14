package main

import (
	"fmt"
	"log"
	"menu/command"
	"menu/files"
	"os/exec"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func main() {
	app := app.New()

	mainWindow := app.NewWindow("MENU")
	ic, err := fyne.LoadResourceFromPath("./logo.png")
	if err != nil {
		log.Println(err.Error())
	}
	mainWindow.SetIcon(ic)

	separ := layout.NewSpacer()

	vault := command.NewVault(files.NewJsonDb("data.json"))

	open := false

	buttons := container.NewVBox()

	add_but := widget.NewButton("Add Script", drawForm(app, mainWindow, buttons, separ, open, vault))

	exit_but := widget.NewButton("Exit", func() { mainWindow.Close() })

	drawButtons(mainWindow, buttons, vault)

	mainWindow.SetContent(container.NewVBox(add_but, buttons, separ, exit_but))

	mainWindow.CenterOnScreen()
	mainWindow.Resize(fyne.NewSize(300, 400))
	mainWindow.Show()

	app.Run()
}

func drawButtons(w fyne.Window, buttons *fyne.Container, vault *command.VaultWithDb) {
	buttons.RemoveAll()

	for _, val := range vault.Commands {
		del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() { deleteCommand(w, val.Id, buttons, vault) })
		ex := widget.NewButton(val.Title, func() { execCommand(w, val.Cmd) })

		buttons.Add(container.New(layout.NewFormLayout(), del, ex))
	}

}

func drawForm(app fyne.App, w fyne.Window, buttons *fyne.Container, separ fyne.CanvasObject, open bool, vault *command.VaultWithDb) func() {
	return func() {
		if !open {
			w2 := app.NewWindow("ADD SCRIPT FORM")

			title_label := widget.NewLabel("Title")
			title_input := widget.NewEntry()
			cmd_label := widget.NewLabel("Command")
			cmd_input := widget.NewEntry()
			empty := widget.NewLabel("")

			save := widget.NewButton("Save", func() {
				addCommand(w2, title_input.Text, cmd_input.Text, vault)
			})

			cancel := widget.NewButton("Cancel", func() {
				w2.Close()
			})

			grid := container.New(layout.NewFormLayout(), title_label, title_input, cmd_label, cmd_input, separ, empty, separ, save, separ, cancel)
			w2.SetContent(grid)
			w2.CenterOnScreen()
			w2.Resize(fyne.NewSize(400, 200))
			w2.Show()
			open = true
			w2.SetOnClosed(func() {
				open = false
				drawButtons(w, buttons, vault)
			})
		}
	}
}

func addCommand(w fyne.Window, title_text string, text_cmd string, vault *command.VaultWithDb) {
	id := uuid.New().String()
	arr := strings.Split(strings.TrimSpace(text_cmd), " ")
	com, err := command.NewComand(id, title_text, arr)
	if err != nil {
		dialog.ShowError(err, w)
		return
	}
	vault.AddCommand(*com)
	w.Close()
}

func deleteCommand(w fyne.Window, id string, buttons *fyne.Container, vault *command.VaultWithDb) {
	dialog.ShowConfirm("DELETE", "Confirm deletion", func(b bool) {
		if b {
			isDelete := vault.DeleteCommadById(id)
			if isDelete {
				drawButtons(w, buttons, vault)
			} else {
				fmt.Println("Ненайдено")
			}
		}

	}, w)

}

func execCommand(w fyne.Window, cmd []string) {
	executer := exec.Command(cmd[0], cmd[1:]...)
	err := executer.Run()
	if err != nil {
		log.Println(err.Error())
		dialog.ShowError(err, w)
	}
}
