// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/maruel/anim1d"
	"github.com/maruel/halloween2018/common"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/apa102"
	"periph.io/x/periph/host"
)

// https://godoc.org/github.com/maruel/anim1d

const pattern = `
{
  "Left": {
    "Curve": "ease-out",
    "Patterns": [
      "#000f00",
      "#00ff00",
      "#1f0f00",
      "#ffa900"
    ],
    "ShowMS": 100,
    "TransitionMS": 700,
    "_type": "Loop"
  },
  "Offset": "50%",
  "Right": {
    "Curve": "ease-out",
    "Patterns": [
      "#1f0f00",
      "#ffa900",
      "#000f00",
      "#00ff00"
    ],
    "ShowMS": 100,
    "TransitionMS": 700,
    "_type": "Loop"
  },
  "_type": "Split"
}
`

const patternAlt = `"#1a5e1f"`

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	s, err := spireg.Open("")
	if err != nil {
		return err
	}
	defer s.Close()
	if err := s.LimitSpeed(2 * physic.MegaHertz); err != nil {
		return err
	}

	opts := apa102.DefaultOpts
	opts.NumPixels = 150
	opts.Temperature = 6500
	a, err := apa102.New(s, &opts)
	if err != nil {
		return err
	}

	var pat anim1d.SPattern
	if err := json.Unmarshal([]byte(pattern), &pat); err != nil {
		return err
	}
	var alt anim1d.SPattern
	if err := json.Unmarshal([]byte(patternAlt), &alt); err != nil {
		return err
	}

	motion := make(chan struct{})
	// TODO(maruel): Get motion from MQTT.
	println("coffin")
	go common.RunDisplay(a, pat.Pattern, alt.Pattern, motion)
	die := make(chan struct{})
	<-die
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "coffin: %s.\n", err)
		os.Exit(1)
	}
}
