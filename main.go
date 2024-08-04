package main

import (
	"encoding/json"
	"fmt"
	"log"
	"menu/commands"
	"menu/files"
	"os/exec"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/widget"
)

func main() {
	a := app.New()
	w := a.NewWindow("MENU")

	s := layout.NewSpacer()

	open := false

	obj, err := read()
	if err != nil {
		log.Println(err.Error())
	}

	add := widget.NewButton("Add Script", func() {
		if !open {
			w2 := a.NewWindow("FORM")
			w2.SetContent(widget.NewLabel("More content"))

			label1 := widget.NewLabel("Title")
			value1 := widget.NewEntry()
			label2 := widget.NewLabel("Command")
			value2 := widget.NewEntry()
			label3 := widget.NewLabel("")
			save := widget.NewButton("Save", func() { log.Println("Save") })
			cancel := widget.NewButton("Cancel", func() { log.Println("Cancel") })
			grid := container.New(layout.NewFormLayout(), label1, value1, label2, value2, s, label3, s, save, s, cancel)
			w2.SetContent(grid)
			w2.CenterOnScreen()
			w2.Resize(fyne.NewSize(400, 200))
			w2.Show()
			open = true
			w2.SetOnClosed(func() {
				log.Println("close")
				open = false
			})
		}

	})

	var buttonts []fyne.CanvasObject
	buttonts = append(buttonts, add)
	for _, val := range obj {
		log.Println(val)
		buttonts = append(buttonts, widget.NewButton(val.Title, execCommand(val.Cmd)))
	}

	exit := widget.NewButton("Exit", func() { w.Close() })
	buttonts = append(buttonts, s, exit)

	w.SetContent(container.NewVBox(buttonts...))

	w.CenterOnScreen()
	w.Resize(fyne.NewSize(300, 400))
	w.Show()

	a.Run()
	tidyUp()
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
