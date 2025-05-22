// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package homin

import (
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo/dapp"
	"github.com/osyah/go-pletyvo/delivery"
	"github.com/osyah/hryzun/buffer"
)

const (
	ChannelTypePrivate = 0 + iota
	ChannelTypePublic
)

type ChannelItem struct {
	Key   uuid.UUID
	Value string
	Hash  dapp.Hash
}

type LocalChannel struct {
	*delivery.Channel
	Type    uint8
	Content *buffer.Ring[*ChannelItem]
}

func (lc LocalChannel) Title() string {
	if lc.Type == ChannelTypePublic {
		return lc.Name + " 💬"
	}

	return lc.Name
}

func (lc LocalChannel) Description() string {
	last, ok := lc.Content.Last()
	if ok {
		return last.Value
	}

	return "Go to the channel to get the latest message!"
}

func (lc LocalChannel) FilterValue() string { return lc.Name }
