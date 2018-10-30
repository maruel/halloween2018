// Copyright 2018 Marc-Antoine Ruel. All rights reserved.
// Use of this source code is governed under the Apache License, Version 2.0
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"testing"

	"github.com/maruel/anim1d"
)

func TestPattern(t *testing.T) {
	var pat anim1d.SPattern
	if err := json.Unmarshal([]byte(pattern), &pat); err != nil {
		t.Fatal(err)
	}
	var alt anim1d.SPattern
	if err := json.Unmarshal([]byte(patternAlt), &alt); err != nil {
		t.Fatal(err)
	}
}
