package commands

import "errors"

type Command struct {
	Id    string   `json:"id"`
	Title string   `json:"title"`
	Cmd   []string `json:"cmd"`
}

func NewComand(id string, title string, cmd []string) (*Command, error) {
	if id == "" {
		return nil, errors.New("ID_EMPTY_STING")
	}
	if title == "" {
		return nil, errors.New("TITLE_EMPTY_STING")
	}
	if len(cmd) == 0 {
		return nil, errors.New("TITLE_EMPTY_STING")
	}

	newCommnd := &Command{Id: id, Title: title, Cmd: cmd}

	return newCommnd, nil
}
