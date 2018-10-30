// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package common

import (
	"image"
	"log"
	"runtime"
	"time"

	"github.com/maruel/anim1d"
	"periph.io/x/periph/conn/display"
	"periph.io/x/periph/conn/gpio"
	"periph.io/x/periph/conn/physic"
)

func RunPIR(pir gpio.PinIn, motion chan<- struct{}) {
	n := pir.Read()
	for {
		if pir.WaitForEdge(-1) {
			if l := pir.Read(); l != n {
				n = l
				if n {
					motion <- struct{}{}
				}
			}
		}
	}
}

func RunDisplay(d display.Drawer, p, alt anim1d.Pattern, motion <-chan struct{}) {
	fps := 30 * physic.Hertz
	if runtime.NumCPU() > 1 {
		fps = 60 * physic.Hertz
	}
	cur := p
	b := d.Bounds()
	f := make(anim1d.Frame, b.Dx())
	t := time.NewTicker(fps.Duration())
	defer t.Stop()
	start := time.Now()
	for now := start; ; {
		since := uint32(now.Sub(start) / time.Millisecond)
		select {
		case _, ok := <-motion:
			if !ok {
				log.Printf("motion channel closed")
				return
			}
			log.Printf("Motion!")
			cur = &anim1d.Transition{
				Before:       anim1d.SPattern{alt},
				After:        anim1d.SPattern{p},
				OffsetMS:     2500 + since,
				TransitionMS: 1000,
				Curve:        anim1d.EaseOut,
			}
		default:
		}
		cur.Render(f, since)
		if err := d.Draw(b, f, image.Point{}); err != nil {
			log.Fatalf("Draw failed: %v", err)
		}
		if tr, ok := cur.(*anim1d.Transition); ok {
			if tr.OffsetMS+tr.TransitionMS < since {
				cur = tr.After.Pattern
				start = start.Add(time.Duration(tr.OffsetMS) * time.Millisecond)
			}
		}
		now = <-t.C
	}
}
