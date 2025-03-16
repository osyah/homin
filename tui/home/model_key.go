// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package home

import "github.com/charmbracelet/bubbles/key"

type ModelKeyMap struct {
	Join   key.Binding
	Create key.Binding
	Leave  key.Binding
}

var ModelKeys = ModelKeyMap{
	Join: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "join channel"),
	),
	Create: key.NewBinding(
		key.WithKeys("ctrl+n"),
		key.WithHelp("ctrl+n", "create channel"),
	),
	Leave: key.NewBinding(
		key.WithKeys("ctrl+l"),
		key.WithHelp("ctrl+l", "leave channel"),
	),
}
