// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package tui

import (
	"log"
	"os"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/osyah/go-pletyvo/client/adapter/dapphttp"
	"github.com/osyah/go-pletyvo/client/adapter/deliveryhttp"
	"github.com/osyah/go-pletyvo/client/engine/http"
	"github.com/osyah/go-pletyvo/protocol/dapp/crypto"
	"github.com/osyah/go-pletyvo/protocol/delivery"

	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
	"github.com/osyah/homin/service"
)

type App struct {
	ctx     *context.Context
	client  *delivery.Service
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

			app.ctx.Signer = crypto.NewED25519(key.PrivateKey[:32])
		} else {
			app.ctx.Page = context.LoginPage
		}
	}

	engine := http.New(http.Config{URL: app.ctx.Config.Gateway})
	app.client = deliveryhttp.New(
		engine, app.ctx.Signer, dapphttp.NewEvent(engine),
	)

	app.program = tea.NewProgram(
		NewModel(app.ctx, service.New(app.client)),
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
