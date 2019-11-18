// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// Package driver contains all code for controlling devices
package driver

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

var supportCoolingProfilesTest = []struct {
	in  []int
	out bool
}{
	{[]int{2, 9, 9}, false},
	{[]int{2, 0, 0}, false},
	{[]int{3, 0, 0}, true},
	{[]int{6, 0, 0}, true},
	{[]int{6, 0, 2}, true},
}

func TestSupportsCoolingProfiles(t *testing.T) {
	for _, tt := range supportCoolingProfilesTest {
		t.Run(fmt.Sprintf("%d.%d.%d", tt.in[0], tt.in[1], tt.in[2]), func(t *testing.T) {
			kraken := KrakenDriver{FirmwareVersion: tt.in, CoolingProfiles: true}
			assert.Equal(t, tt.out, kraken.SupportsCoolingProfiles())
		})
	}
}
