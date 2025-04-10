// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/charmbracelet/bubbles/list"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
)

type Contact struct {
	locals map[string]config.Contact
}

func NewContact() *Contact {
	return &Contact{locals: make(map[string]config.Contact)}
}

func (c *Contact) Get(*context.Context) ([]list.Item, error) {
	var err error

	if len(c.locals) == 0 {
		c.locals, err = config.GetContacts()
		if err != nil {
			return nil, err
		}
	}

	contacts := make([]list.Item, len(c.locals))

	var i uint16

	for hash, local := range c.locals {
		contacts[i] = &homin.LocalContact{
			Name: local.Name, Hash: hash,
		}

		i++
	}

	return contacts, nil
}

func (c Contact) Delete(contact *homin.LocalContact) error {
	delete(c.locals, contact.Hash)

	return config.SaveContacts(c.locals)
}
