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
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
	"github.com/google/uuid"
)

func main() {
	app := app.New()
	mainWindow := app.NewWindow("MENU")

	separ := layout.NewSpacer()

	open := false

	buttons := container.NewVBox()

	add := widget.NewButton("Add Script", drawForm(app, buttons, separ, open))

	exit := widget.NewButton("Exit", func() { mainWindow.Close() })

	drawButtons(buttons)

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

func drawButtons(buttons *fyne.Container) {
	buttons.RemoveAll()
	obj, err := read()
	if err != nil {
		log.Println(err.Error())
	}

	for _, val := range obj {
		log.Println(val)

		buttons.Add(container.New(layout.NewFormLayout(), widget.NewButton("X", func() { deleteCommand(val.Id, buttons) }), widget.NewButton(val.Title, execCommand(val.Cmd))))
	}
}

func deleteCommand(id string, buttons *fyne.Container) {
	log.Println(id)
	var result []commands.Command
	coms, err := read()
	if err != nil {
		log.Println(err.Error())
	}
	for i, val := range coms {
		if val.Id == id {
			result = append(coms[:i], coms[i+1:]...)
		}
	}
	l, err := json.Marshal(result)
	if err != nil {
		log.Println(err.Error())
	}
	log.Println(coms)

	files.WriteFile(l, "data.json")
	drawButtons(buttons)
}

func drawForm(app fyne.App, buttons *fyne.Container, separ fyne.CanvasObject, open bool) func() {
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
				com, _ := commands.NewComand(id, title_input.Text, arr)
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
				drawButtons(buttons)
			})
		}
	}
}

func execCommand(cmd []string) func() {
	return func() {
		cmd := exec.Command(cmd[0], cmd[1:]...)
		err := cmd.Run()
		if err != nil {
			log.Fatal(err)
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
		log.Println("!!", err.Error())
		return nil, err
	}
	err = json.Unmarshal(file, &coms)
	if err != nil {
		log.Println("!!", err.Error())
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
