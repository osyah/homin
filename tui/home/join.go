// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package home

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type JoinModel struct {
	ctx     *context.Context
	service *service.Home

	textInput textinput.Model
}

func NewJoinModel(ctx *context.Context, service *service.Home) JoinModel {
	textInput := textinput.New()
	textInput.Placeholder = "channel/00000000-0000-0000-0000-000000000000"
	textInput.ShowSuggestions = true
	textInput.SetSuggestions([]string{"channel", "chat"})
	textInput.Focus()

	return JoinModel{ctx: ctx, service: service, textInput: textInput}
}

func (JoinModel) Init() tea.Cmd { return nil }

func (jm JoinModel) Resize(wsm *tea.WindowSizeMsg) JoinModel {
	jm.textInput.Width = wsm.Width

	return jm
}

func (jm JoinModel) Update(msg tea.Msg) (JoinModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			jm.ctx.Page = context.HomePage
		case tea.KeyEnter:
			local, err := jm.service.Join(jm.textInput.Value())
			if err != nil {
				break
			}

			jm.textInput.SetValue("")

			jm.ctx.Page = context.HomePage

			return jm, func() tea.Msg { return context.JoinChannel{Local: local} }
		}
	}

	var cmd tea.Cmd

	jm.textInput, cmd = jm.textInput.Update(msg)

	return jm, cmd
}

func (jm JoinModel) View() string {
	builder := strings.Builder{}

	builder.WriteString(
		"To join a channel or chat, you need to enter the link.\n" +
			"The link format should be 'type/id', where type contains channel or chat.\n\n",
	)
	builder.WriteString(jm.textInput.View())

	return builder.String()
}
