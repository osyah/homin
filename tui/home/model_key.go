// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package home

import "github.com/charmbracelet/bubbles/key"

type ModelKeyMap struct {
	Join key.Binding
}

var ModelKeys = ModelKeyMap{
	Join: key.NewBinding(
		key.WithKeys("ctrl+j"),
		key.WithHelp("ctrl+j", "join channel"),
	),
}
