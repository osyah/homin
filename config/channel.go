// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"os"

	"github.com/google/uuid"

	"github.com/osyah/homin"
)

type Channel struct {
	Type uint8     `json:"type"`
	ID   uuid.UUID `json:"id"`
}

func SaveChannels(channels []Channel) error {
	f, err := os.Create(homin.Path + "/channels.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(channels, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func GetChannels() ([]Channel, error) {
	b, err := os.ReadFile(homin.Path + "/channels.json")
	if err != nil {
		return nil, err
	}

	var channels []Channel

	if err = json.Unmarshal(b, &channels); err != nil {
		return nil, err
	}

	return channels, nil
}
