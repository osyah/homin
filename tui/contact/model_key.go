// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package contact

import "github.com/charmbracelet/bubbles/key"

type ModelKeyMap struct {
	Delete key.Binding
}

var ModelKeys = ModelKeyMap{
	Delete: key.NewBinding(
		key.WithKeys("ctrl+d"),
		key.WithHelp("ctrl+d", "delete contact"),
	),
}
