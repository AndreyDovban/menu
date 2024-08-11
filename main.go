package main

import (
	"encoding/json"
	"fmt"
	"log"
	"menu/commands"
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

	separ := layout.NewSpacer()

	open := false

	buttons := container.NewVBox()

	add := widget.NewButton("Add Script", drawForm(app, mainWindow, buttons, separ, open))

	exit := widget.NewButton("Exit", func() { mainWindow.Close() })

	drawButtons(mainWindow, buttons)

	mainWindow.SetContent(container.NewVBox(add, buttons, separ, exit))

	mainWindow.CenterOnScreen()
	mainWindow.Resize(fyne.NewSize(300, 400))
	mainWindow.Show()

	app.Run()
	tidyUp()
}

func KeyDown(ev *fyne.KeyEvent) {
	log.Println(ev.Name)
}

func drawButtons(w fyne.Window, buttons *fyne.Container) {
	buttons.RemoveAll()
	obj, err := read()
	if err != nil {
		log.Println(err.Error())
		dialog.ShowError(err, w)
	}

	for _, val := range obj {
		log.Println(val)

		del := widget.NewButtonWithIcon("", theme.DeleteIcon(), func() { deleteCommand(w, val.Id, buttons) })
		ex := widget.NewButton(val.Title, execCommand(w, val.Cmd))

		buttons.Add(container.New(layout.NewFormLayout(), del, ex))
	}
}

func deleteCommand(w fyne.Window, id string, buttons *fyne.Container) {
	dialog.ShowConfirm("DELETE", "Confirm deletion", func(b bool) {
		if b {
			log.Println(id)
			var result []commands.Command
			coms, err := read()
			if err != nil {
				log.Println(err.Error())
				dialog.ShowError(err, w)
			}
			for i, val := range coms {
				if val.Id == id {
					result = append(coms[:i], coms[i+1:]...)
				}
			}
			l, err := json.Marshal(result)
			if err != nil {
				log.Println(err.Error())
				dialog.ShowError(err, w)
			}
			log.Println(coms)

			files.WriteFile(l, "data.json")
			drawButtons(w, buttons)
		}

	}, w)

}

func drawForm(app fyne.App, w fyne.Window, buttons *fyne.Container, separ fyne.CanvasObject, open bool) func() {
	return func() {
		if !open {
			w2 := app.NewWindow("ADD SCRIPT FORM")
			w2.SetContent(widget.NewLabel("More content"))

			title_label := widget.NewLabel("Title")
			title_input := widget.NewEntry()
			cmd_label := widget.NewLabel("Command")
			cmd_input := widget.NewEntry()
			empty := widget.NewLabel("")
			save := widget.NewButton("Save", func() {
				log.Println("Save")
				id := uuid.New().String()
				log.Println("&&&", id)
				arr := strings.Split(strings.TrimSpace(cmd_input.Text), " ")
				com, err := commands.NewComand(id, title_input.Text, arr)
				if err != nil {
					dialog.ShowError(err, w2)
					return
				}
				write(*com)
				w2.Close()
			})
			cancel := widget.NewButton("Cancel", func() {
				log.Println("Cancel")
				w2.Close()
			})
			grid := container.New(layout.NewFormLayout(), title_label, title_input, cmd_label, cmd_input, separ, empty, separ, save, separ, cancel)
			w2.SetContent(grid)
			w2.CenterOnScreen()
			w2.Resize(fyne.NewSize(400, 200))
			w2.Show()
			open = true
			w2.SetOnClosed(func() {
				log.Println("close")
				open = false
				drawButtons(w, buttons)
			})
		}
	}
}

func execCommand(w fyne.Window, cmd []string) func() {
	return func() {
		cmd := exec.Command(cmd[0], cmd[1:]...)
		err := cmd.Run()
		if err != nil {
			log.Println(err.Error())
			dialog.ShowError(err, w)
		}
	}
}

func tidyUp() {
	fmt.Println("Exited")
}

func read() ([]commands.Command, error) {
	var coms []commands.Command
	file, err := files.ReadFile("./data.json")
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	err = json.Unmarshal(file, &coms)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return coms, nil
}

func write(com commands.Command) {

	coms, err := read()
	if err != nil {
		log.Println(err.Error())
	}
	coms = append(coms, com)
	l, err := json.Marshal(coms)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(coms)

	files.WriteFile(l, "data.json")
}
