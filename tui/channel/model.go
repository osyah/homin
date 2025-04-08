// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package channel

import (
	"strings"

	"github.com/charmbracelet/bubbles/textarea"
	"github.com/charmbracelet/bubbles/viewport"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/protocol/delivery"

	"github.com/osyah/homin"
	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type Model struct {
	ctx     *context.Context
	service *service.Channel

	viewPort viewport.Model
	textArea textarea.Model
}

func NewModel(ctx *context.Context, service *service.Channel) Model {
	viewPort := viewport.New(100, 10)

	textArea := textarea.New()
	textArea.SetHeight(3)
	textArea.Placeholder = "Write a message..."
	textArea.ShowLineNumbers = false
	textArea.KeyMap.InsertNewline.SetEnabled(false)
	textArea.Focus()

	return Model{
		ctx:      ctx,
		service:  service,
		viewPort: viewPort,
		textArea: textArea,
	}
}

func (Model) Init() tea.Cmd { return nil }

func (m Model) Resize(wsm *tea.WindowSizeMsg) Model {
	m.viewPort.Width = wsm.Width
	m.viewPort.Height = wsm.Height - m.textArea.Height() - 2

	m.textArea.SetWidth(wsm.Width)

	if m.ctx.Channel != nil {
		m.viewPort.SetContent(
			m.renderContent(m.ctx.Channel.Content.Get()),
		)
		m.viewPort.GotoBottom()
	}

	return m
}

func (m Model) Update(msg tea.Msg) (Model, tea.Cmd) {
	switch msg := msg.(type) {
	case context.UpdateContent:
		m = m.updateChannelContent()
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyEsc:
			m.ctx.Page = context.HomePage
			m.viewPort.SetContent("Loading...")
		case tea.KeyEnter:
			switch m.ctx.Channel.Type {
			case homin.ChannelTypePrivate:
				input := &delivery.PostCreateInput{
					PostInput: &delivery.PostInput{
						Channel: m.ctx.Channel.Hash,
						Content: m.textArea.Value(),
					},
				}

				item, err := m.service.CreatePost(m.ctx, input)
				if err != nil {
					break
				}

				m.ctx.Channel.Content.Add(item)
				m.viewPort.SetContent(
					m.renderContent(m.ctx.Channel.Content.Get()),
				)

				m.textArea.Reset()
				m.viewPort.GotoBottom()
			case homin.ChannelTypePublic:
				input := &delivery.MessageInput{
					ID:      uuid.Must(uuid.NewV7()),
					Channel: m.ctx.Channel.Hash,
					Content: m.textArea.Value(),
				}

				item, err := m.service.SendMessage(m.ctx, input)
				if err != nil {
					break
				}

				m.ctx.Channel.Content.Add(item)
				m.viewPort.SetContent(
					m.renderContent(m.ctx.Channel.Content.Get()),
				)

				m.textArea.Reset()
				m.viewPort.GotoBottom()
			}
		}
	}

	var vpCmd, taCmd tea.Cmd

	m.viewPort, vpCmd = m.viewPort.Update(msg)
	m.textArea, taCmd = m.textArea.Update(msg)

	return m, tea.Batch(vpCmd, taCmd)
}

func (m Model) View() string {
	builder := strings.Builder{}

	builder.WriteString("--- ")
	builder.WriteString(m.ctx.Channel.Title())
	builder.WriteString(" ---\n")

	builder.WriteString(m.viewPort.View())
	builder.WriteString("\n\n")

	if m.ctx.Channel.Type == homin.ChannelTypePublic {
		builder.WriteString(m.textArea.View())
	} else {
		if m.ctx.Channel.Author == m.ctx.Signer.Address() {
			builder.WriteString(m.textArea.View())
		}
	}

	return builder.String()
}
