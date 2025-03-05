// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package tui

import (
	"time"

	"github.com/osyah/homin/context"
)

func (a App) ListenEvents() {
	for {
		time.Sleep(time.Second * 3)

		if a.ctx.Channel == nil {
			continue
		}

		a.program.Send(context.UpdateContent{})
	}
}
