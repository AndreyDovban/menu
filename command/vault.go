package command

import (
	"encoding/json"
	"fmt"
	"menu/files"

	"time"
)

type Vault struct {
	Commands []Command `json:"commands"`
	UpdateAt time.Time `json:"updateAt"`
}

type VaultWithDb struct {
	Vault
	db files.JsonDb
}

func NewVault(db *files.JsonDb) *VaultWithDb {
	data, err := db.Read()
	if err != nil {
		return &VaultWithDb{
			Vault: Vault{
				Commands: []Command{},
				UpdateAt: time.Now(),
			},
			db: *db,
		}
	}

	var vault Vault

	err = json.Unmarshal(data, &vault)
	if err != nil {
		fmt.Println(err.Error())
		return &VaultWithDb{
			Vault: Vault{
				Commands: []Command{},
				UpdateAt: time.Now(),
			},
			db: *db,
		}
	}

	return &VaultWithDb{
		Vault: vault,
		db:    *db,
	}
}

func (vdb *VaultWithDb) AddCommand(cmd Command) {
	vdb.Vault.Commands = append(vdb.Vault.Commands, cmd)
	vdb.save()
}

func (vdb *VaultWithDb) DeleteCommadById(id string) bool {
	var commands []Command
	isDelete := false
	for _, com := range vdb.Vault.Commands {
		if com.Id != id {
			commands = append(commands, com)
			continue
		}
		isDelete = true
	}

	vdb.Vault.Commands = commands
	vdb.save()

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

func (vdb *VaultWithDb) save() {
	vdb.Vault.UpdateAt = time.Now()
	data, err := vdb.Vault.ToBytes()
	if err != nil {
		fmt.Println(err.Error())
	}
	vdb.db.Write(data)
}
