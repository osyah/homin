// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package homin

import (
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/protocol/delivery"
	"github.com/osyah/hryzun/buffer"
)

const (
	ChannelTypePrivate = 0 + iota
	ChannelTypePublic
)

type ChannelItem struct {
	Key   uuid.UUID
	Value string
}

type LocalChannel struct {
	*delivery.Channel
	Type uint8

	Posts    *buffer.Ring[*ChannelItem]
	Messages *buffer.Ring[*ChannelItem]
}

func (lc LocalChannel) Title() string {
	if lc.Type == ChannelTypePublic {
		return lc.Name + " 💬"
	}

	return lc.Name
}

func (lc LocalChannel) Description() string {
	switch lc.Type {
	case ChannelTypePrivate:
		last, ok := lc.Posts.Last()
		if ok {
			return last.Value
		}
	case ChannelTypePublic:
		last, ok := lc.Messages.Last()
		if ok {
			return last.Value
		}
	}

	return "Go to the channel to get the latest message!"
}

func (lc LocalChannel) FilterValue() string { return lc.Name }
