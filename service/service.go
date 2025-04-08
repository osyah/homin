// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/osyah/go-pletyvo/client/adapter/dapphttp"
	"github.com/osyah/go-pletyvo/client/adapter/deliveryhttp"
	"github.com/osyah/go-pletyvo/client/engine/http"

	"github.com/osyah/homin/context"
)

type Service struct {
	Login   *Login
	Home    *Home
	Channel *Channel
	Contact *Contact
}

func New(ctx *context.Context) *Service {
	engine := http.New(http.Config{URL: ctx.Config.Gateway})

	eventService := dapphttp.NewEvent(engine)
	deliveryClient := deliveryhttp.New(engine, ctx.Signer, eventService)

	contact := NewContact()

	return &Service{
		Login:   NewLogin(),
		Home:    NewHome(deliveryClient.Channel, eventService),
		Channel: NewChannel(deliveryClient, eventService, contact),
		Contact: contact,
	}
}
