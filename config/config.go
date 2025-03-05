// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"os"

	"github.com/osyah/homin"
)

type Config struct {
	Gateway string `json:"gateway"`
	Auth    *Auth  `json:"auth,omitempty"`
}

func (c Config) Save() error {
	f, err := os.Create(homin.Path + "/config.json")
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(&c, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func New() (*Config, error) {
	var config Config

	b, err := os.ReadFile(homin.Path + "/config.json")
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(b, &config)
	if err != nil {
		return nil, err
	}

	return &config, nil
}
