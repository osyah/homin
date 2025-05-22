// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/osyah/go-pletyvo"
	"github.com/osyah/go-pletyvo/dapp"
	"github.com/osyah/go-pletyvo/delivery"

	"github.com/osyah/homin/context"
)

type Service struct {
	Login   *Login
	Home    *Home
	Channel *Channel
	Contact *Contact
}

func New(ctx *context.Context) *Service {
	engine := pletyvo.NewEngine(pletyvo.EngineConfig{URL: ctx.Config.Gateway})

	eventService := dapp.NewEventClient(engine)
	deliveryClient := delivery.NewClient(engine, ctx.Signer, eventService)

	contact := NewContact()

	return &Service{
		Login:   NewLogin(),
		Home:    NewHome(deliveryClient.Channel, eventService),
		Channel: NewChannel(deliveryClient, eventService, contact),
		Contact: contact,
	}
}
