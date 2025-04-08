// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"os"

	"github.com/osyah/homin"
)

var DefaultContacts = map[string]Contact{
	"u4OwwMAMQSfUnNQGJePONaAVVipaOACFz0vDfyzIUKQ": {
		Name: "Осяг",
	},
}

type Contact struct {
	Name string `json:"name"`
}

func SaveContacts(contacts map[string]Contact) error {
	f, err := os.Create(homin.Path + "/contacts.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(contacts, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func GetContacts() (map[string]Contact, error) {
	b, err := os.ReadFile(homin.Path + "/contacts.json")
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultContacts, SaveContacts(DefaultContacts)
		}

		return nil, err
	}

	var contacts map[string]Contact

	if err = json.Unmarshal(b, &contacts); err != nil {
		return nil, err
	}

	return contacts, nil
}
