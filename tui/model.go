// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package tui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
	"github.com/osyah/homin/tui/auth"
	"github.com/osyah/homin/tui/channel"
	"github.com/osyah/homin/tui/home"
)

type Model struct {
	ctx *context.Context

	loginModel   auth.LoginModel
	homeModel    home.Model
	channelModel channel.Model
	joinModel    home.JoinModel
	createModel  home.CreateModel
}

func NewModel(ctx *context.Context, service *service.Service) Model {
	return Model{
		ctx:          ctx,
		loginModel:   auth.NewLoginModel(ctx, service.Login),
		homeModel:    home.NewModel(ctx, service.Home),
		channelModel: channel.NewModel(ctx, service.Channel),
		joinModel:    home.NewJoinModel(ctx, service.Home),
		createModel:  home.NewCreateModel(ctx, service.Home),
	}
}

func (Model) Init() tea.Cmd { return nil }

func (m Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.channelModel = m.channelModel.Resize(&msg)
		m.homeModel = m.homeModel.Resize(&msg)
		m.joinModel = m.joinModel.Resize(&msg)
		m.createModel = m.createModel.Resize(&msg)
	case tea.KeyMsg:
		switch msg.Type {
		case tea.KeyCtrlC:
			return m, tea.Quit
		}
	}

	var cmd tea.Cmd

	switch m.ctx.Page {
	case context.LoginPage:
		m.loginModel, cmd = m.loginModel.Update(msg)
	case context.HomePage:
		m.homeModel, cmd = m.homeModel.Update(msg)
	case context.ChannelPage:
		m.channelModel, cmd = m.channelModel.Update(msg)
	case context.JoinPage:
		m.joinModel, cmd = m.joinModel.Update(msg)
	case context.CreatePage:
		m.createModel, cmd = m.createModel.Update(msg)
	}

	return m, cmd
}

func (m Model) View() string {
	switch m.ctx.Page {
	case context.LoginPage:
		return m.loginModel.View()
	case context.HomePage:
		return m.homeModel.View()
	case context.ChannelPage:
		return m.channelModel.View()
	case context.JoinPage:
		return m.joinModel.View()
	case context.CreatePage:
		return m.createModel.View()
	default:
		return m.homeModel.View()
	}
}
