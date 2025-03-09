// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package home

import (
	"log"

	"github.com/charmbracelet/bubbles/key"
	"github.com/charmbracelet/bubbles/list"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/osyah/homin"
	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

var homeModelStyle = lipgloss.NewStyle().Padding(1, 2)

type Model struct {
	ctx  *context.Context
	list list.Model
}

func NewModel(ctx *context.Context, service *service.Home) Model {
	channels, err := service.GetChannels()
	if err != nil {
		log.Fatalln("homin/tui/home: " + err.Error())
	}

	delegate := list.NewDefaultDelegate()
	delegate.UpdateFunc = delegateUpdate

	list := list.New(channels, delegate, 20, 20)
	list.Title = "Channels"
	list.AdditionalFullHelpKeys = func() []key.Binding {
		return []key.Binding{ModelKeys.Join}
	}

	return Model{ctx: ctx, list: list}
}

func delegateUpdate(msg tea.Msg, model *list.Model) tea.Cmd {
	switch msg := msg.(type) {
	case context.JoinChannel:
		return model.InsertItem((len(model.Items()) + 1), msg.Local)
	}

	return nil
}

func (Model) Init() tea.Cmd { return nil }

func (m Model) Resize(wsm *tea.WindowSizeMsg) Model {
	h, v := homeModelStyle.GetFrameSize()
	m.list.SetSize((wsm.Width - h), (wsm.Height - v))

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return m, tea.Quit
		case tea.KeyCtrlJ:
			m.ctx.Page = context.JoinPage
		case tea.KeyEnter:
			m.ctx.Channel = m.list.SelectedItem().(*homin.LocalChannel)
			m.ctx.Page = context.ChannelPage

			return m, func() tea.Msg { return context.UpdateContent{} }
		}
	}

	var cmd tea.Cmd

	m.list, cmd = m.list.Update(msg)

	return m, cmd
}

func (m Model) View() string {
	return homeModelStyle.Render(m.list.View())
}
