// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"context"

	"github.com/charmbracelet/bubbles/list"
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/protocol/delivery"
	"github.com/osyah/hryzun/buffer"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
)

type Home struct {
	client   delivery.ChannelService
	locals   []config.Channel
	channels map[uuid.UUID]*delivery.Channel
}

func NewHome(client delivery.ChannelService) *Home {
	return &Home{
		client:   client,
		channels: make(map[uuid.UUID]*delivery.Channel),
	}
}

func (h *Home) GetChannels() ([]list.Item, error) {
	var err error

	if h.locals == nil {
		h.locals, err = config.GetChannels()
		if err != nil {
			return nil, err
		}
	}

	channels := make([]list.Item, len(h.locals))

	for i, local := range h.locals {
		channel, ok := h.channels[local.ID]
		if !ok {
			channel, err = h.client.GetByID(context.Background(), local.ID)
			if err != nil {
				return nil, err
			}

			h.channels[channel.ID] = channel
		}

		channels[i] = &homin.LocalChannel{
			Channel:  channel,
			Type:     local.Type,
			Posts:    buffer.NewRing[*homin.ChannelItem](20),
			Messages: buffer.NewRing[*homin.ChannelItem](30),
		}
	}

	return channels, nil
}
