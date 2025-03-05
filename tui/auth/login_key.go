// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package auth

import "github.com/charmbracelet/bubbles/key"

type LoginKeyMap struct {
	Quit     key.Binding
	Generate key.Binding
	Continue key.Binding
}

func (lkm LoginKeyMap) ShortHelp() []key.Binding {
	return []key.Binding{lkm.Quit, lkm.Generate, lkm.Continue}
}

func (lkm LoginKeyMap) FullHelp() [][]key.Binding {
	return [][]key.Binding{{lkm.Quit, lkm.Generate, lkm.Continue}}
}

var LoginKeys = LoginKeyMap{
	Quit: key.NewBinding(
		key.WithKeys("esc", "ctrl+c"),
		key.WithHelp("ctrl+c", "quit"),
	),
	Generate: key.NewBinding(
		key.WithKeys("ctrl+g"),
		key.WithHelp("ctrl+g", "generate"),
	),
	Continue: key.NewBinding(
		key.WithKeys("enter"),
		key.WithHelp("enter", "continue"),
	),
}
