// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package tui

import (
	"fmt"
	"os"

	"github.com/google/uuid"

	"github.com/osyah/homin"
	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
)

func MakeConfig(ctx *context.Context) error {
	err := os.MkdirAll(homin.Path, os.ModePerm)
	if err != nil {
		return err
	}

	var channels []config.Channel

	switch homin.Mode {
	case homin.ModeLocal:
		ctx.Config = &config.Config{
			Gateway: "http://localhost:8049/api",
		}
	case homin.ModeTest:
		ctx.Config = &config.Config{
			Gateway: "http://testnet.pletyvo.osyah.com/api",
		}

		channels = []config.Channel{
			{Type: 0, ID: uuid.MustParse("0195672b-634f-7077-aeb2-d7c658a8d08d")},
			{Type: 1, ID: uuid.MustParse("0195672b-634f-7077-aeb2-d7c658a8d08d")},
		}
	default:
		return fmt.Errorf("homin/tui: invalid mode")
	}

	err = ctx.Config.Save()
	if err != nil {
		return err
	}

	return config.SaveChannels(channels)
}
