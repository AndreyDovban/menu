package command

import (
	"encoding/json"
	"fmt"
	"menu/files"
	"time"
)

type Vault struct {
	Commands []Command `json:"acccounts"`
	UpdateAt time.Time `json:"updateAt"`
}

func NewVault() *Vault {
	data, err := files.ReadFile("data.json")
	if err != nil {
		return &Vault{
			Commands: []Command{},
			UpdateAt: time.Now(),
		}
	}

	var vault Vault

	err = json.Unmarshal(data, &vault)
	if err != nil {
		fmt.Println(err.Error())
		return &Vault{
			Commands: []Command{},
			UpdateAt: time.Now(),
		}
	}

	return &vault
}

func (vault *Vault) AddCommand(cmd Command) {
	vault.Commands = append(vault.Commands, cmd)
	vault.save()
}

func (vault *Vault) DeleteCommadById(id string) bool {
	var commands []Command
	isDelete := false
	for _, com := range vault.Commands {
		if com.Id != id {
			commands = append(commands, com)
			continue
		}
		isDelete = true
	}

	vault.Commands = commands
	vault.save()

	return isDelete
}

func (vault *Vault) ToBytes() ([]byte, error) {
	file, err := json.Marshal(vault)
	if err != nil {
		fmt.Println(err.Error())
		return nil, err
	}

	return file, nil
}

func (vault *Vault) save() {
	vault.UpdateAt = time.Now()
	data, err := vault.ToBytes()
	if err != nil {
		fmt.Println(err.Error())
	}
	files.WriteFile(data, "data.json")
}
