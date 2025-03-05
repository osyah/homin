// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package main

import (
	"log"

	"github.com/osyah/homin/tui"
)

func main() {
	app, err := tui.NewApp()
	if err != nil {
		log.Fatalln("homin/tui: " + err.Error())
	}

	app.Run()
}
