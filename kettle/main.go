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
  "Curve": "ease-out",
  "Patterns": [
    {
      "Patterns": [
        {
          "Child": {
            "Child": {
              "Curve": "direct",
              "Left": "#ffa900",
              "Right": "#000000",
              "_type": "Gradient"
            },
            "Length": 16,
            "Offset": 0,
            "_type": "Subset"
          },
          "MovePerHour": 108000,
          "_type": "Rotate"
        },
        {
          "_type": "Aurore"
        }
      ],
      "_type": "Add"
    },
    {
      "Patterns": [
        {
          "_type": "Aurore"
        },
        {
          "C": "#ffffff",
          "_type": "NightStars"
        }
      ],
      "_type": "Add"
    }
  ],
  "ShowMS": 10000,
  "TransitionMS": 5000,
  "_type": "Loop"
}
`

const patternAlt = `
{
  "Curve": "ease-out",
  "Patterns": [
    "#ff0000",
    "#0f0000"
  ],
  "ShowMS": 10,
  "TransitionMS": 100,
  "_type": "Loop"
}
`

const patternIncoming = `
{
  "Curve": "ease-out",
  "Patterns": [
    {
      "Patterns": [
        {
          "Child": "Lff0000ff0000ee0000dd0000cc0000bb0000aa0000990000880000770000660000550000440000330000220000110000",
          "MovePerHour": 108000,
          "_type": "Rotate"
        },
        {
          "_type": "Aurore"
        }
      ],
      "_type": "Add"
    },
    {
      "Patterns": [
        {
          "_type": "Aurore"
        },
        {
          "C": "#ffffff",
          "_type": "NightStars"
        }
      ],
      "_type": "Add"
    }
  ],
  "ShowMS": 10000,
  "TransitionMS": 5000,
  "_type": "Loop"
}
`

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
	opts.NumPixels = 45
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
	println("kettle")
	go common.RunDisplay(a, pat.Pattern, alt.Pattern, motion)
	die := make(chan struct{})
	<-die
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "kettle: %s.\n", err)
		os.Exit(1)
	}
}
