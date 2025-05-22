// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package tui

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/osyah/go-pletyvo/dapp"

	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type App struct {
	ctx     *context.Context
	program *tea.Program
}

func NewApp() (*App, error) {
	app := &App{
		ctx: &context.Context{Page: context.HomePage},
	}

	var err error

	app.ctx.Config, err = config.New()
	if err != nil {
		if os.IsNotExist(err) {
			if err := MakeConfig(app.ctx); err != nil {
				return nil, err
			}

			app.ctx.Page = context.LoginPage
		} else {
			return nil, err
		}
	} else {
		if app.ctx.Config.Auth != nil {
			key, err := config.GetKey(app.ctx.Config.Auth.Key)
			if err != nil {
				return nil, err
			}

			app.ctx.Signer = dapp.NewED25519(key.PrivateKey[:32])
		} else {
			app.ctx.Page = context.LoginPage
		}
	}

	app.program = tea.NewProgram(
		NewModel(app.ctx, service.New(app.ctx)),
		tea.WithAltScreen(),
	)

	return app, nil
}

func (a App) Run() {
	go a.ListenEvents()

	if _, err := a.program.Run(); err != nil {
		log.Fatalln("homin/tui: " + err.Error())
	}
}
