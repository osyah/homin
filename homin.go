// Copyright (c) 2025 Osyah
// SPDX-License-Identifier: MIT

package homin

import (
	"flag"
	"log"
	"os"
	"path/filepath"
	"strings"
)

const (
	ModeLocal = iota
	ModeTest
)

var (
	Mode = func() uint8 {
		var mode string

		flag.StringVar(&mode, "mode", "main", "Select application mode.")
		flag.Parse()

		switch strings.ToLower(mode) {
		case "local":
			return ModeLocal
		case "test":
			return ModeTest
		default:
			return ModeTest
		}
	}()

	Path = func() string {
		flag.Parse()
		value := os.Getenv("HOMIN_PATH")
		if value != "" {
			if filepath.IsAbs(value) {
				log.Fatalln("homin: invalid '$HOMIN_PATH' value")
			}

			return value
		}

		dir, err := os.UserConfigDir()
		if err != nil {
			log.Fatalln("homin: " + err.Error())
		}

		switch Mode {
		case ModeLocal:
			return dir + "/Osyah/Homin Local"
		case ModeTest:
			return dir + "/Osyah/Homin Test"
		default:
			return dir + "/Osyah/Homin"
		}
	}()
)
