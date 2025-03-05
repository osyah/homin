// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package auth

import (
	"strings"

	"github.com/charmbracelet/bubbles/help"
	"github.com/charmbracelet/bubbles/textarea"
	tea "github.com/charmbracelet/bubbletea"

	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type LoginModel struct {
	ctx     *context.Context
	service *service.Login

	textArea textarea.Model
	help     help.Model
}

func NewLoginModel(ctx *context.Context, service *service.Login) LoginModel {
	textArea := textarea.New()
	textArea.Placeholder = "If it is missing, generate it using ctrl+g"
	textArea.ShowLineNumbers = false
	textArea.Focus()

	return LoginModel{
		ctx:      ctx,
		service:  service,
		textArea: textArea,
		help:     help.New(),
	}
}

func (lm LoginModel) Init() tea.Cmd { return nil }

func (lm LoginModel) Update(msg tea.Msg) (LoginModel, tea.Cmd) {
	var cmd tea.Cmd

	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		if msg.Width < 100 {
			lm.textArea.SetWidth(msg.Width)
		} else {
			lm.textArea.SetWidth(msg.Width / 2)
		}

		lm.help.Width = msg.Width
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			return lm, tea.Quit
		case tea.KeyEnter:
			if err := lm.service.SaveKey(lm.ctx, lm.textArea.Value(), ""); err != nil {
				lm.textArea.SetValue("An error occurred while logging in: " + err.Error())

				break
			}

			lm.ctx.Page = context.HomePage
		case tea.KeyCtrlG:
			mnemonic, err := lm.service.GenerateMnemonic()
			if err != nil {
				lm.textArea.SetValue("An error occurred during generation: " + err.Error())

				break
			}

			lm.textArea.SetValue(mnemonic)
		}
	}

	lm.textArea, cmd = lm.textArea.Update(msg)

	return lm, cmd
}

func (lm LoginModel) View() string {
	builder := strings.Builder{}

	builder.WriteString("Welcome, new ruler of the world!\n")
	builder.WriteString("Before you can use the app, you have to enter a mnemonic phrase.\n\n")

	builder.WriteString(lm.textArea.View())

	builder.WriteString("\n\n")
	builder.WriteString(lm.help.View(LoginKeys))

	return builder.String()
}
