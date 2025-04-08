// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package contact

import (
	"log"

	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

var contactModelStyle = lipgloss.NewStyle().Padding(1, 2)

type Model struct {
	ctx     *context.Context
	service *service.Contact

	list list.Model
}

func NewModel(ctx *context.Context, service *service.Contact) Model {
	items, err := service.Get(ctx)
	if err != nil {
		log.Fatalln("homin/tui/contact: " + err.Error())
	}

	list := list.New(items, list.NewDefaultDelegate(), 20, 20)
	list.Title = "Contacts"

	return Model{ctx: ctx, service: service, list: list}
}

func (Model) Init() tea.Cmd { return nil }

func (m Model) Resize(wsm *tea.WindowSizeMsg) Model {
	h, v := contactModelStyle.GetFrameSize()
	m.list.SetSize((wsm.Width - h), (wsm.Height - v))

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			if !m.list.FilteringEnabled() {
				return m, tea.Quit
			}
		}
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return contactModelStyle.Render(m.list.View())
}
