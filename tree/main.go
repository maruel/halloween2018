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
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
	"periph.io/x/periph/conn/spi/spireg"
	"periph.io/x/periph/devices/apa102"
	"periph.io/x/periph/host"
	"periph.io/x/periph/host/bcm283x"
)

// https://godoc.org/github.com/maruel/anim1d

const pattern = `
{
	"_type": "Add",
	"Patterns": [
		{
				"C": "#ff9000",
				"_type": "NightStars"
		},
		{
			"Curve": "steps(1,end)",
			"Patterns": [
				"#0f0f0f",
				"#000000",
				"#0f0f0f",
				"#000000",
				"#000000",
				"#0f0f0f",
				"#000000",
				"#000000",
				"#000000",
				"#000000"
			],
			"ShowMS": 100,
			"TransitionMS": 0,
			"_type": "Loop"
		}
	]
}
`

const patternAlt = `
{
	"Curve": "steps(1,end)",
	"Patterns": ["#8f8f8f","#000000","#8f8f8f","#000000","#000000","#8f8f8f","#000000","#000000","#000000","#000000"],
	"ShowMS": 100,
	"TransitionMS": 0,
	"_type": "Loop"
}
`

const pattern2 = `
		{
			"Curve": "steps(1,end)",
			"Patterns": [
				"#0f0f0f",
				"#000000",
				"#0f0f0f",
				"#000000",
				"#000000",
				"#0f0f0f",
				"#000000",
				"#000000",
				"#000000",
				"#000000"
			],
			"ShowMS": 100,
			"TransitionMS": 0,
			"_type": "Loop"
		}
`

const pattern3 = `
{
	"_type": "Rotate",
	"Child": {
		"_type": "Dim",
		"Child": "Rainbow",
		"Intensity": 20
	},
	"MovePerHour": 7200
}
`

func mainImpl() error {
	if _, err := host.Init(); err != nil {
		return err
	}

	pir := bcm283x.GPIO4
	if err := pir.In(gpio.Float, gpio.BothEdges); err != nil {
		return err
	}
	defer pir.Halt()

	s, err := spireg.Open("")
	if err != nil {
		return err
	}
	defer s.Close()
	if err := s.LimitSpeed(2 * physic.MegaHertz); err != nil {
		return err
	}

	opts := apa102.DefaultOpts
	opts.NumPixels = 300
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
	println("tree")
	go common.RunDisplay(a, pat.Pattern, alt.Pattern, motion)
	go common.RunPIR(pir, motion)
	die := make(chan struct{})
	<-die
	return nil
}

func main() {
	if err := mainImpl(); err != nil {
		fmt.Fprintf(os.Stderr, "tree: %s.\n", err)
		os.Exit(1)
	}
}
