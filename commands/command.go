package commands

import "errors"

type Command struct {
	Title string   `json:"title"`
	Cmd   []string `json:"cmd"`
}

func NewComand(title string, cmd []string) (*Command, error) {
	if title == "" {
		return nil, errors.New("TITLE_EMPTY_STING")
	}
	// if len(cmd == 0) {
	// 	return nil, errors.New("TITLE_EMPTY_STING")
	// }

	newCommnd := &Command{Title: title, Cmd: cmd}

	return newCommnd, nil
}
