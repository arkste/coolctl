// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// Package driver contains all code for controlling devices
package driver

import (
	"image/color"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestColorFromHexString(t *testing.T) {
	colorCode, err := colorFromHexString("ff0000")

	assert.Equal(t, &color.RGBA{255, 0, 0, 1}, colorCode)
	assert.Nil(t, err)
}

func TestColorFromHexStringInvalid(t *testing.T) {
	colorCode, err := colorFromHexString("foobar")

	assert.Nil(t, colorCode)
	assert.Error(t, err)
}

func TestPaletteFromColors(t *testing.T) {
	palette, err := paletteFromColors([]string{"FF0000", "00FF00", "0000FF"})

	assert.Equal(t, &color.Palette{
		&color.RGBA{255, 0, 0, 1},
		&color.RGBA{0, 255, 0, 1},
		&color.RGBA{0, 0, 255, 1},
	}, palette)

	assert.Nil(t, err)
}

func TestPaletteFromColorsInvalid(t *testing.T) {
	palette, err := paletteFromColors([]string{"foobar"})

	assert.Nil(t, palette)
	assert.Error(t, err)
}

func TestMakeRange(t *testing.T) {
	tmpRange := makeRange(10, 22, 2)

	assert.Len(t, tmpRange, 6)
	assert.Equal(t, []int{10, 12, 14, 16, 18, 20}, tmpRange)
}

func TestParseProfile(t *testing.T) {
	profile := parseProfile("20 25  35 25  50 55  60 100")

	assert.Equal(t, SpeedProfile{{20, 25}, {35, 25}, {50, 55}, {60, 100}}, profile)
}

var normalizeTests = []struct {
	in  string
	out SpeedProfile
}{
	{"35 25  20 25  50 55  50 100", SpeedProfile{{20, 25}, {35, 25}, {50, 55}, {60, 100}}},
	{"35 25  20 25  50 55  60 90", SpeedProfile{{20, 25}, {35, 25}, {50, 55}, {60, 100}}},
	{"35 25  20 25  50 55  50 90", SpeedProfile{{20, 25}, {35, 25}, {50, 55}, {50, 90}, {60, 100}}},
}

func TestNormalizeProfile(t *testing.T) {
	for _, tnt := range normalizeTests {
		t.Run(tnt.in, func(t *testing.T) {
			profile := parseProfile(tnt.in)
			assert.Equal(t, tnt.out, normalizeProfile(profile, criticalTemp))
		})
	}
}

func TestInterpolateProfile(t *testing.T) {
	profile := parseProfile("20 25  35 25  50 55  60 100")
	assert.Equal(t, SpeedProfile{{20, 25}, {22, 25}, {24, 25}, {26, 25}, {28, 25}, {30, 25}, {32, 25}, {34, 25}, {36, 27}, {38, 31}, {40, 35}, {42, 39}, {44, 43}, {46, 47}, {48, 51}, {50, 55}, {52, 64}, {54, 73}, {56, 82}, {58, 91}, {60, 100}}, interpolateProfile(profile))
}

func TestNormalizeInterpolateProfile(t *testing.T) {
	profile := normalizeProfile(parseProfile("20 25  35 25  50 55  60 100"), criticalTemp)
	assert.Equal(t, SpeedProfile{{20, 25}, {22, 25}, {24, 25}, {26, 25}, {28, 25}, {30, 25}, {32, 25}, {34, 25}, {36, 27}, {38, 31}, {40, 35}, {42, 39}, {44, 43}, {46, 47}, {48, 51}, {50, 55}, {52, 64}, {54, 73}, {56, 82}, {58, 91}, {60, 100}}, interpolateProfile(profile))
}
