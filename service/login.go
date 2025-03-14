// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package service

import (
	"github.com/osyah/go-pletyvo/protocol/dapp/crypto"
	"github.com/tyler-smith/go-bip39"

	"github.com/osyah/homin/config"
	"github.com/osyah/homin/context"
)

type Login struct{}

func NewLogin() *Login { return &Login{} }

func (Login) GenerateMnemonic() (string, error) {
	entropy, err := bip39.NewEntropy(256)
	if err != nil {
		return "", err
	}

	return bip39.NewMnemonic(entropy)
}

func (Login) SaveKey(ctx *context.Context, phrase, password string) error {
	if !bip39.IsMnemonicValid(phrase) {
		return bip39.ErrInvalidMnemonic
	}

	ctx.Signer = crypto.NewED25519(bip39.NewSeed(phrase, password)[:32])

	if err := config.SaveKey(ctx.Signer.Address(), config.Key{
		PrivateKey: ctx.Signer.(*crypto.ED25519).Private(),
	}); err != nil {
		return err
	}

	ctx.Config.Auth = &config.Auth{Key: ctx.Signer.Address()}

	return ctx.Config.Save()
}
