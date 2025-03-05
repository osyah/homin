// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package config

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/osyah/go-pletyvo/protocol/dapp"

	"github.com/osyah/homin"
)

type Auth struct {
	Key dapp.Hash `json:"key"`
}

type Key struct {
	PrivateKey []byte `json:"private_key"`
}

func SaveKey(address dapp.Hash, key Key) error {
	if err := os.MkdirAll((homin.Path + "/keys"), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(fmt.Sprintf(("%s/keys/%s.json"), homin.Path, address))
	if err != nil {
		return err
	}
	defer f.Close()

	b, err := json.MarshalIndent(&key, "", "  ")
	if err != nil {
		return err
	}

	_, err = f.Write(b)

	return err
}

func GetKey(address dapp.Hash) (*Key, error) {
	path := fmt.Sprintf("%s/keys/%s.json", homin.Path, address)

	b, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var key Key

	err = json.Unmarshal(b, &key)
	if err != nil {
		return nil, err
	}

	return &key, nil
}
