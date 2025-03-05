// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import "github.com/osyah/go-pletyvo/protocol/delivery"

type Service struct {
	Login   *Login
	Home    *Home
	Channel *Channel
}

func New(service *delivery.Service) *Service {
	return &Service{
		Login:   NewLogin(),
		Home:    NewHome(service.Channel),
		Channel: NewChannel(service),
	}
}
