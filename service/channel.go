// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"strings"
	"time"

	"github.com/osyah/go-pletyvo"
	"github.com/osyah/go-pletyvo/dapp"
	"github.com/osyah/go-pletyvo/delivery"

	"github.com/osyah/homin"
	"github.com/osyah/homin/context"
)

type Channel struct {
	message delivery.MessageService
	post    delivery.PostService
	event   dapp.EventService
	contact *Contact
}

func NewChannel(client *delivery.Service, event dapp.EventService, contact *Contact) *Channel {
	return &Channel{
		message: client.Message,
		post:    client.Post,
		event:   event,
		contact: contact,
	}
}

func (c Channel) GetPosts(ctx *context.Context, option *pletyvo.QueryOption) ([]*delivery.Post, error) {
	return c.post.Get(ctx.Background(), ctx.Channel.ID, option)
}

func (c Channel) GetMessages(ctx *context.Context, option *pletyvo.QueryOption) ([]*dapp.Event, error) {
	return c.message.Get(ctx.Background(), ctx.Channel.ID, option)
}

func (c Channel) FormatPost(post *delivery.Post) (*homin.ChannelItem, error) {
	post.Content = delivery.PrepareContent(post.Content)
	if len(post.Content) == 0 {
		return nil, delivery.ErrEmptyContent
	}

	return c.renderPost(post), nil
}

func (c Channel) renderPost(post *delivery.Post) *homin.ChannelItem {
	var builder strings.Builder

	builder.WriteString("\033[32m")
	builder.WriteString(time.Unix(post.ID.Time().UnixTime()).Format("02/01 15:04 "))
	builder.WriteString("\033[0m")

	builder.WriteString(post.Content)

	return &homin.ChannelItem{
		Key:   post.ID,
		Value: builder.String(),
		Hash:  post.Hash,
	}
}

func (c Channel) FormatMessage(event *dapp.Event) (*homin.ChannelItem, error) {
	var input delivery.MessageInput

	if err := event.Body.Unmarshal(&input); err != nil {
		return nil, err
	}

	input.Content = delivery.PrepareContent(input.Content)
	if len(input.Content) == 0 {
		return nil, delivery.ErrEmptyContent
	}

	return c.renderMessage(event, &input), nil
}

func (c Channel) renderMessage(event *dapp.Event, input *delivery.MessageInput) *homin.ChannelItem {
	author := dapp.NewHash(event.Auth.Schema, event.Auth.PublicKey).String()

	var builder strings.Builder

	builder.WriteString("\033[32m")
	builder.WriteString(time.Unix(input.ID.Time().UnixTime()).Format("02/01 15:04 "))
	builder.WriteString("\033[0m")

	contact, ok := c.contact.locals[author]
	if ok {
		builder.WriteString("\033[1;33m")
		builder.WriteString(contact.Name)
		builder.WriteString("\033[0m")
	} else {
		builder.WriteString("\033[90m")
		builder.WriteString(author[:5])
		builder.WriteString("...")
		builder.WriteString(author[38:])
		builder.WriteString("\033[0m")
	}

	builder.WriteByte(' ')
	builder.WriteString(input.Content)

	return &homin.ChannelItem{
		Key:   input.ID,
		Value: builder.String(),
		Hash:  dapp.NewHash(event.Auth.Schema, event.Auth.Signature),
	}
}

func (c Channel) CreatePost(ctx *context.Context, input *delivery.PostCreateInput) (*homin.ChannelItem, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	input.Content = delivery.PrepareContent(input.Content)
	if len(input.Content) == 0 {
		return nil, delivery.ErrEmptyContent
	}

	var body dapp.EventBody

	last, ok := ctx.Channel.Content.Last()
	if ok {
		body = dapp.NewEventBody(
			dapp.EventBodyLinked, dapp.JSONDataType, delivery.PostCreate, input,
		)
		body.SetParent(last.Hash)
	} else {
		body = dapp.NewEventBody(
			dapp.EventBodyBasic, dapp.JSONDataType, delivery.PostCreate, input,
		)
	}

	event := &dapp.EventInput{Body: body, Auth: ctx.Signer.Auth(body)}

	response, err := c.event.Create(ctx.Background(), event)
	if err != nil {
		return nil, err
	}

	return c.renderPost(&delivery.Post{
		ID:      response.ID,
		Author:  ctx.Signer.Address(),
		Hash:    dapp.NewHash(event.Auth.Schema, event.Auth.Signature),
		Channel: ctx.Channel.ID,
		Content: input.Content,
	}), nil
}

func (c Channel) SendMessage(ctx *context.Context, input *delivery.MessageInput) (*homin.ChannelItem, error) {
	err := input.Validate()
	if err != nil {
		return nil, err
	}

	input.Content = delivery.PrepareContent(input.Content)
	if len(input.Content) == 0 {
		return nil, delivery.ErrEmptyContent
	}

	body := dapp.NewEventBody(
		dapp.EventBodyBasic, dapp.JSONDataType, delivery.MessageCreate, &input,
	)

	event := &dapp.EventInput{Body: body, Auth: ctx.Signer.Auth(body)}

	if err = c.message.Send(ctx.Background(), event); err != nil {
		return nil, err
	}

	return c.renderMessage(&dapp.Event{ID: input.ID, Body: event.Body, Auth: event.Auth}, input), nil
}
