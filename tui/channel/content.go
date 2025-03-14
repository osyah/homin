// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package channel

import (
	"slices"
	"strings"

	"github.com/osyah/go-pletyvo"
	"github.com/osyah/go-pletyvo/protocol/delivery"

	"github.com/osyah/homin"
)

func (m Model) updateChannelContent() Model {
	var err error

	switch m.ctx.Channel.Type {
	case homin.ChannelTypePrivate:
		var posts []*delivery.Post

		last, ok := m.ctx.Channel.Content.Last()
		if !ok {
			posts, err = m.service.GetPosts(m.ctx, &pletyvo.QueryOption{Limit: 20})
			if err != nil {
				if err == pletyvo.CodeNotFound {
					m.viewPort.SetContent(
						"Unfortunately, there are no posts on this channel.",
					)
				}

				break
			}

			slices.Reverse(posts)
		} else {
			posts, err = m.service.GetPosts(
				m.ctx, &pletyvo.QueryOption{After: last.Key, Order: true},
			)
			if err != nil {
				if err != pletyvo.CodeNotFound {
					break
				}
			}
		}

		for _, post := range posts {
			if post == nil {
				continue
			}

			m.ctx.Channel.Content.Add(m.service.FormatPost(post))
		}

		m.viewPort.SetContent(m.renderContent(m.ctx.Channel.Content.Get()))
	case homin.ChannelTypePublic:
		var messages []*delivery.Message

		last, ok := m.ctx.Channel.Content.Last()
		if !ok {
			messages, err = m.service.GetMessages(m.ctx, &pletyvo.QueryOption{Limit: 20})
			if err != nil {
				if err == pletyvo.CodeNotFound {
					m.viewPort.SetContent(
						"This channel is empty, be the first to make history!",
					)
				}

				break
			}

			slices.Reverse(messages)
		} else {
			messages, err = m.service.GetMessages(
				m.ctx, &pletyvo.QueryOption{After: last.Key, Order: true},
			)
			if err != nil {
				if err != pletyvo.CodeNotFound {
					break
				}
			}
		}

		for _, message := range messages {
			if message == nil {
				continue
			}

			m.ctx.Channel.Content.Add(m.service.FormatMessage(message))
		}

		m.viewPort.SetContent(m.renderContent(m.ctx.Channel.Content.Get()))
	}

	m.viewPort.GotoBottom()

	return m
}

func (m Model) renderContent(items []*homin.ChannelItem) string {
	switch len(items) {
	case 0:
		return ""
	case 1:
		return items[0].Value
	}

	builder := strings.Builder{}

	builder.WriteString(items[0].Value)

	for _, item := range items[1:] {
		if item == nil {
			continue
		}

		builder.WriteString("\n")
		builder.WriteString(item.Value)
	}

	return builder.String()
}
