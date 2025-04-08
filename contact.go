// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package homin

type LocalContact struct {
	Name, Hash string
}

func (lc LocalContact) Title() string { return lc.Name }

func (lc LocalContact) Description() string { return lc.Hash }

func (lc LocalContact) FilterValue() string { return lc.Name }
