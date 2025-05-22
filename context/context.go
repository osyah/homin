// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package context

import (
	"context"

	"github.com/osyah/go-pletyvo/dapp"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
)

const (
	LoginPage = 1 + iota
	HomePage
	ChannelPage
	JoinPage
	CreatePage
	ContactPage
)

type Context struct {
	Page    uint8
	Config  *config.Config
	Signer  dapp.Signer
	Channel *homin.LocalChannel
}

func (Context) Background() context.Context { return context.Background() }

type UpdateContent struct{}

type JoinChannel struct{ Local *homin.LocalChannel }
