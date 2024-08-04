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

			title_label := widget.NewLabel("Title")
			title_input := widget.NewEntry()
			cmd_label := widget.NewLabel("Command")
			cmd_input := widget.NewEntry()
			empty := widget.NewLabel("")
			save := widget.NewButton("Save", func() {
				log.Println("Save")
				arr := []string{"firefox"}
				com, _ := commands.NewComand("TITLE", arr)
				u := []commands.Command{*com}
				write(u)
				w2.Close()
			})
			cancel := widget.NewButton("Cancel", func() {
				log.Println("Cancel")
				w2.Close()
			})
			grid := container.New(layout.NewFormLayout(), title_label, title_input, cmd_label, cmd_input, s, empty, s, save, s, cancel)
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

func write(com []commands.Command) {
	log.Println(com)
	l, err := json.Marshal(com)
	if err != nil {
		log.Println(err.Error())
	}
	files.WriteFile(l, "data.json")

}
