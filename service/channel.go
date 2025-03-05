// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/google/uuid"
	"github.com/osyah/go-pletyvo"
	"github.com/osyah/go-pletyvo/protocol/dapp"
	"github.com/osyah/go-pletyvo/protocol/dapp/crypto"
	"github.com/osyah/go-pletyvo/protocol/delivery"

	"github.com/osyah/homin"
	"github.com/osyah/homin/context"
)

var (
	timeStyle   = lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	senderStyle = lipgloss.NewStyle().Foreground(lipgloss.Color("8"))
)

type Channel struct {
	message delivery.MessageService
	post    delivery.PostService
}

func NewChannel(service *delivery.Service) *Channel {
	return &Channel{message: service.Message, post: service.Post}
}

func (c Channel) GetPosts(ctx *context.Context, option *pletyvo.QueryOption) ([]*delivery.Post, error) {
	return c.post.Get(ctx.Background(), ctx.Channel.ID, option)
}

func (c Channel) GetMessages(ctx *context.Context, option *pletyvo.QueryOption) ([]*delivery.Message, error) {
	return c.message.Get(ctx.Background(), ctx.Channel.ID, option)
}

func (c Channel) FormatPost(post *delivery.Post) *homin.ChannelItem {
	var builder strings.Builder

	builder.WriteString(
		timeStyle.Render(
			time.Unix(post.ID.Time().UnixTime()).Format("02/01 15:04 "),
		),
	)
	builder.WriteString(post.Content)

	return &homin.ChannelItem{
		Key:   post.ID,
		Value: builder.String(),
	}
}

func (c Channel) FormatMessage(message *delivery.Message) *homin.ChannelItem {
	var (
		builder strings.Builder
		input   delivery.MessageInput
		author  dapp.Hash
	)

	if err := message.Body.Unmarshal(&input); err != nil {
		return &homin.ChannelItem{
			Key:   uuid.Must(uuid.NewV7()),
			Value: "<invalid message>",
		}
	}

	author = crypto.NewHash(message.Auth.Schema, message.Auth.PublicKey)

	builder.WriteString(
		timeStyle.Render(
			time.Unix(input.ID.Time().UnixTime()).Format("02/01 15:04 "),
		),
	)
	builder.WriteString(
		senderStyle.Render(author.String()[:5] + "..." + author.String()[38:] + " "),
	)
	builder.WriteString(input.Content)

	return &homin.ChannelItem{
		Key:   input.ID,
		Value: builder.String(),
	}
}

func (c Channel) CreatePost(ctx *context.Context, input *delivery.PostCreateInput) (*delivery.Post, error) {
	response, err := c.post.Create(ctx.Background(), input)
	if err != nil {
		return nil, err
	}

	return &delivery.Post{
		ID:      response.ID,
		Content: input.Content,
	}, nil
}

func (c Channel) SendMessage(ctx *context.Context, input *delivery.MessageInput) (*delivery.Message, error) {
	if err := input.Validate(); err != nil {
		return nil, err
	}

	body := dapp.NewEventBody(
		dapp.EventBodyBasic, dapp.JSONDataType, delivery.MessageCreate, &input,
	)

	message := &delivery.Message{Body: body, Auth: ctx.Signer.Auth(body)}

	return message, c.message.Send(ctx.Background(), message)
}
