// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package home

import (
	"strings"

	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/osyah/go-pletyvo/protocol/delivery"

	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type CreateModel struct {
	ctx     *context.Context
	service *service.Home

	textInput textinput.Model
}

func NewCreateModel(ctx *context.Context, service *service.Home) CreateModel {
	textInput := textinput.New()
	textInput.Placeholder = "channel name..."
	textInput.Focus()

	return CreateModel{ctx: ctx, service: service, textInput: textInput}
}

func (CreateModel) Init() tea.Cmd { return nil }

func (cm CreateModel) Resize(wsm *tea.WindowSizeMsg) CreateModel {
	cm.textInput.Width = wsm.Width

	return cm
}

func (cm CreateModel) Update(msg tea.Msg) (CreateModel, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			cm.ctx.Page = context.HomePage
		case tea.KeyEnter:
			input := &delivery.ChannelCreateInput{
				ChannelInput: &delivery.ChannelInput{
					Name: cm.textInput.Value(),
				},
			}

			locals, err := cm.service.Create(cm.ctx, input)
			if err != nil {
				break
			}

			cm.textInput.SetValue("")
			cm.ctx.Page = context.HomePage

			return cm, tea.Batch(
				func() tea.Msg { return context.JoinChannel{Local: locals[0]} },
				func() tea.Msg { return context.JoinChannel{Local: locals[1]} },
			)
		}
	}

	var cmd tea.Cmd

	cm.textInput, cmd = cm.textInput.Update(msg)

	return cm, cmd
}

func (cm CreateModel) View() string {
	builder := strings.Builder{}

	builder.WriteString(
		"Before you create a channel, you'll have to come up with a name.\n\n",
	)
	builder.WriteString(cm.textInput.View())

	return builder.String()
}
