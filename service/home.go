// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/bubbles/list"
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/protocol/dapp"
	"github.com/osyah/go-pletyvo/protocol/delivery"
	"github.com/osyah/hryzun/buffer"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
)

var StringToChannelType = map[string]uint8{
	"channel": homin.ChannelTypePrivate,
	"chat":    homin.ChannelTypePublic,
}

type Home struct {
	client   delivery.ChannelService
	event    dapp.EventService
	locals   []config.Channel
	channels map[uuid.UUID]*delivery.Channel
}

func NewHome(client delivery.ChannelService, event dapp.EventService) *Home {
	return &Home{
		client:   client,
		event:    event,
		channels: make(map[uuid.UUID]*delivery.Channel),
	}
}

func (h *Home) GetChannels(ctx *context.Context) ([]list.Item, error) {
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
			channel, err = h.client.GetByID(ctx.Background(), local.ID)
			if err != nil {
				return nil, err
			}

			h.channels[channel.ID] = channel
		}

		channels[i] = &homin.LocalChannel{
			Channel: channel,
			Type:    local.Type,
			Content: buffer.NewRing[*homin.ChannelItem](ctx.Config.BufferSize),
		}
	}

	return channels, nil
}

func (h *Home) Join(ctx *context.Context, link string) (*homin.LocalChannel, error) {
	v := strings.Split(link, "/")

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
		channel, err = h.client.GetByID(ctx.Background(), id)
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
		Channel: channel,
		Type:    channelType,
		Content: buffer.NewRing[*homin.ChannelItem](30),
	}, nil
}

func (h *Home) Create(ctx *context.Context, input *delivery.ChannelCreateInput) ([2]*homin.LocalChannel, error) {
	if err := input.Validate(); err != nil {
		return [2]*homin.LocalChannel{}, err
	}

	body := dapp.NewEventBody(
		dapp.EventBodyBasic, dapp.JSONDataType,
		delivery.ChannelCreate, input,
	)

	response, err := h.event.Create(ctx.Background(), &dapp.EventInput{
		Body: body, Auth: ctx.Signer.Auth(body),
	})
	if err != nil {
		return [2]*homin.LocalChannel{}, err
	}

	channel, err := h.client.GetByID(ctx.Background(), response.ID)
	if err != nil {
		return [2]*homin.LocalChannel{}, err
	}

	h.channels[response.ID] = channel

	h.locals = append(
		h.locals,
		config.Channel{Type: homin.ChannelTypePrivate, ID: channel.ID},
		config.Channel{Type: homin.ChannelTypePublic, ID: channel.ID},
	)

	if err = config.SaveChannels(h.locals); err != nil {
		return [2]*homin.LocalChannel{}, err
	}

	return [2]*homin.LocalChannel{
		{
			Channel: channel,
			Type:    homin.ChannelTypePrivate,
			Content: buffer.NewRing[*homin.ChannelItem](ctx.Config.BufferSize),
		},
		{
			Channel: channel,
			Type:    homin.ChannelTypePublic,
			Content: buffer.NewRing[*homin.ChannelItem](ctx.Config.BufferSize),
		},
	}, nil
}

func (h *Home) Leave(channel *homin.LocalChannel) error {
	for i, local := range h.locals {
		if local.Type != channel.Type {
			continue
		}

		if local.ID != channel.ID {
			continue
		}

		h.locals = append(h.locals[:i], h.locals[1+i:]...)

		if err := config.SaveChannels(h.locals); err != nil {
			return err
		}

		break
	}

	return nil
}
