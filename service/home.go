// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"context"
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/protocol/delivery"
	"github.com/osyah/hryzun/buffer"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
)

var StringToChannelType = map[string]uint8{
	"channel": homin.ChannelTypePrivate,
	"chat":    homin.ChannelTypePublic,
}

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

func (h *Home) Join(s string) (*homin.LocalChannel, error) {
	v := strings.Split(s, "/")

	if len(v) != 2 {
		return nil, fmt.Errorf("homin/service: invalid join format")
	}

	channelType, ok := StringToChannelType[v[0]]
	if !ok {
		return nil, fmt.Errorf("homin/service: invalid join type")
	}

	id, err := uuid.Parse(v[1])
	if err != nil {
		return nil, err
	}

	channel, ok := h.channels[id]
	if !ok {
		channel, err := h.client.GetByID(context.Background(), id)
		if err != nil {
			return nil, err
		}

		h.channels[channel.ID] = channel
	}

	h.locals = append(h.locals, config.Channel{
		Type: channelType, ID: id,
	})

	if err = config.SaveChannels(h.locals); err != nil {
		return nil, err
	}

	return &homin.LocalChannel{
		Channel:  channel,
		Type:     channelType,
		Posts:    buffer.NewRing[*homin.ChannelItem](20),
		Messages: buffer.NewRing[*homin.ChannelItem](30),
	}, nil
}
