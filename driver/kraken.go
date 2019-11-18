// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// coolctl is a Golang-Port of liquidctl.
// Copyright (C) 2018–2019 Jonas Malaco
// Copyright (C) 2018–2019 each contribution's author

// Package driver contains all code for controlling devices
package driver

import (
	"context"
	"fmt"
	"strconv"
	"time"

	"github.com/google/gousb"
	log "github.com/sirupsen/logrus"
)

const (
	productID = 0x170e // Kraken X (X42, X52, X62 or X72)
	vendorID  = 0x1e71 // NZXT

	readEndpoint  = 1
	readLength    = 64
	writeEndpoint = 1
	writeLength   = 65

	totalLEDs    = 9
	criticalTemp = 60
)

var (
	// Debug is the debug level, 0 = no output, 10 = more output
	Debug int

	config    = 1
	iface     = 0
	alternate = 0
	bufSize   = 0
	timeout   = time.Duration(0)

	speedChannels = map[string][]int{
		"fan":  {0x80, 25, 100},
		"pump": {0xc0, 50, 100},
	}

	colorChannels = map[string]int{
		"sync": 0x0,
		"logo": 0x1,
		"ring": 0x2,
	}

	colorModes = map[string][]int{
		// byte3/mode, byte2/reverse, byte4/modifier, min colors, max colors, only ring (0=no, 1=yes)
		"off":                          {0x00, 0x00, 0x00, 0, 0, 0},
		"fixed":                        {0x00, 0x00, 0x00, 1, 1, 0},
		"super-fixed":                  {0x00, 0x00, 0x00, 1, 9, 0}, // independent logo + ring leds
		"fading":                       {0x01, 0x00, 0x00, 2, 8, 0},
		"spectrum-wave":                {0x02, 0x00, 0x00, 0, 0, 0},
		"backwards-spectrum-wave":      {0x02, 0x10, 0x00, 0, 0, 0},
		"marquee-3":                    {0x03, 0x00, 0x00, 1, 1, 1},
		"marquee-4":                    {0x03, 0x00, 0x08, 1, 1, 1},
		"marquee-5":                    {0x03, 0x00, 0x10, 1, 1, 1},
		"marquee-6":                    {0x03, 0x00, 0x18, 1, 1, 1},
		"backwards-marquee-3":          {0x03, 0x10, 0x00, 1, 1, 1},
		"backwards-marquee-4":          {0x03, 0x10, 0x08, 1, 1, 1},
		"backwards-marquee-5":          {0x03, 0x10, 0x10, 1, 1, 1},
		"backwards-marquee-6":          {0x03, 0x10, 0x18, 1, 1, 1},
		"covering-marquee":             {0x04, 0x00, 0x00, 1, 8, 1},
		"covering-backwards-marquee":   {0x04, 0x10, 0x00, 1, 8, 1},
		"alternating":                  {0x05, 0x00, 0x00, 2, 2, 1},
		"moving-alternating":           {0x05, 0x08, 0x00, 2, 2, 1},
		"backwards-moving-alternating": {0x05, 0x18, 0x00, 2, 2, 1},
		"breathing":                    {0x06, 0x00, 0x00, 1, 8, 0}, // colors for each step
		"super-breathing":              {0x06, 0x00, 0x00, 1, 9, 0}, // one step, independent logo + ring leds
		"pulse":                        {0x07, 0x00, 0x00, 1, 8, 0},
		"tai-chi":                      {0x08, 0x00, 0x00, 2, 2, 1},
		"water-cooler":                 {0x09, 0x00, 0x00, 0, 0, 1},
		"loading":                      {0x0a, 0x00, 0x00, 1, 1, 1},
		"wings":                        {0x0c, 0x00, 0x00, 1, 1, 1},
		"super-wave":                   {0x0d, 0x00, 0x00, 1, 8, 1}, // independent ring leds
		"backwards-super-wave":         {0x0d, 0x10, 0x00, 1, 8, 1}, // independent ring leds
	}

	animationSpeeds = map[string]int{
		"slowest": 0x0,
		"slower":  0x1,
		"normal":  0x2,
		"faster":  0x3,
		"fastest": 0x4,
	}
)

type contextReader interface {
	ReadContext(context.Context, []byte) (int, error)
}

// KrakenDriver holds all driver relevant informations
type KrakenDriver struct {
	ProductID       gousb.ID
	VendorID        gousb.ID
	FirmwareVersion []int
	CoolingProfiles bool
	*gousb.Context
	*gousb.Interface
	*gousb.InEndpoint
	*gousb.OutEndpoint
}

// NewKrakenDriver creates a new USB Context instance & returns a new KrakenDriver
func NewKrakenDriver() *KrakenDriver {
	ctx := gousb.NewContext()
	ctx.Debug(Debug)
	log.SetLevel(log.Level(Debug))

	return &KrakenDriver{
		ProductID: productID,
		VendorID:  vendorID,
		Context:   ctx,
	}
}

// Connect connects to the USB device
func (d *KrakenDriver) Connect() {
	dev, err := d.Context.OpenDeviceWithVIDPID(d.VendorID, d.ProductID)
	if err != nil {
		log.Fatal("NZXT Kraken X (X42, X52, X62 or X72) not found")
	}
	defer dev.Close()

	err = dev.SetAutoDetach(true)
	if err != nil {
		log.Fatal(err)
	}

	cfg, err := dev.Config(config)
	if err != nil {
		log.Fatalf("dev.Config(%d): %v", config, err)
	}

	d.Interface, err = cfg.Interface(iface, alternate)
	if err != nil {
		log.Fatalf("cfg.Interface(%d, %d): %v", iface, alternate, err)
	}

	d.InEndpoint, err = d.Interface.InEndpoint(readEndpoint)
	if err != nil {
		log.Fatalf("dev.InEndpoint(): %s", err)
	}

	d.OutEndpoint, err = d.Interface.OutEndpoint(writeEndpoint)
	if err != nil {
		log.Fatalf("dev.OutEndpoint(): %s", err)
	}
}

// GetStatus reads & returns the current device status
func (d *KrakenDriver) GetStatus() (string, uint64, uint64, string) {
	msg := d.read()

	temperature := fmt.Sprintf("%d.%d", uint64(msg[1]), uint64(msg[2]))
	fanSpeed := uint64(msg[3])<<8 | uint64(msg[4])
	pumpSpeed := uint64(msg[5])<<8 | uint64(msg[6])
	firmwareVersion := d.readFirmwareVersion(msg)

	return temperature, fanSpeed, pumpSpeed, firmwareVersion
}

// SetColor sets the color of a channel & mode
func (d *KrakenDriver) SetColor(channel, mode string, colors []string) {
	colorChannel, ok := colorChannels[channel]
	if !ok {
		log.Fatalf("channel %s not found", channel)
	}

	colorMode, ok := colorModes[mode]
	if !ok {
		log.Fatalf("mode %s not found", mode)
	}

	mval, mod2, mod4, mincolors, maxcolors, ringonly := colorMode[0], colorMode[1], colorMode[2], colorMode[3], colorMode[4], colorMode[5]
	if ringonly == 1 && channel != "ring" {
		log.Fatalf("mode %s unsupported with channel %s", mode, channel)
	}

	palette, err := paletteFromColors(colors)
	if err != nil {
		log.Fatal(err)
	}

	steps := generateSteps(*palette, mincolors, maxcolors, mode, ringonly)
	for seq, step := range steps {
		logoRed, logoGreen, logoBlue, _ := step[0].RGBA()

		buf := []byte{
			0x2,
			0x4c,
			byte(mod2 | colorChannel),
			byte(mval),
			byte(animationSpeeds["normal"] | seq<<5 | mod4),
			byte(logoGreen),
			byte(logoRed),
			byte(logoBlue),
		}

		for _, leds := range step[1:] {
			red, green, blue, _ := leds.RGBA()
			colors := []byte{byte(red), byte(green), byte(blue)}
			buf = append(buf, colors...)
		}

		d.write(buf)
	}
}

// SetSpeed sets a profile for a speed channel
func (d *KrakenDriver) SetSpeed(channel, profile string) {
	speedChannel, ok := speedChannels[channel]
	if !ok {
		log.Fatalf("channel %s not found", channel)
	}

	cbase, dmin, dmax, p := speedChannel[0], speedChannel[1], speedChannel[2], interpolateProfile(normalizeProfile(parseProfile(profile), criticalTemp))
	log.Infof("setting profile for channel '%s': %v", channel, p)

	for i, profile := range p {
		duty := profile[1]

		if duty < dmin {
			duty = dmin
		} else if duty > dmax {
			duty = dmax
		}

		d.write([]byte{0x2, 0x4d, byte(cbase + i), byte(profile[0]), byte(duty)})
	}
}

// SetFixedSpeed checks if device supports cooling profiles and then sets the provided duty for the channel either instant or not
func (d *KrakenDriver) SetFixedSpeed(channel, duty string) {
	if d.SupportsCoolingProfiles() {
		d.SetSpeed(channel, "0 "+duty+"  59 "+duty+"  60 100  100 100")
	} else {
		d.setInstantSpeed(channel, duty)
	}
}

// SupportsCoolingProfiles checks if the current firmware supports cooling profiles
func (d *KrakenDriver) SupportsCoolingProfiles() bool {
	if d.CoolingProfiles == false {
		d.GetStatus()
	}

	return d.FirmwareVersion[0] >= 3 && d.FirmwareVersion[1] >= 0 && d.FirmwareVersion[2] >= 0
}

// read reads from the USB device
func (d *KrakenDriver) read() []byte {
	var rdr contextReader = d.InEndpoint
	if bufSize > 1 {
		log.Print("creating buffer...")
		s, err := d.InEndpoint.NewStream(readLength, bufSize)
		if err != nil {
			log.Fatalf("ep.NewStream(): %v", err)
		}
		defer s.Close()
		rdr = s
	}

	opCtx := context.Background()
	if timeout > 0 {
		var done func()
		opCtx, done = context.WithTimeout(opCtx, timeout)
		defer done()
	}
	msg := make([]byte, readLength)
	_, err := rdr.ReadContext(opCtx, msg)
	if err != nil {
		log.Fatalf("reading from device failed: %v", err)
	}
	log.Infof("reading: %d", msg)
	log.Infof("reading: % 02x", msg)

	return msg
}

// write writes to the USB device
func (d *KrakenDriver) write(data []byte) {
	padding := make([]byte, writeLength-len(data))
	log.Infof("writing: %d", data)
	log.Infof("writing: % 02x", data)
	data = append(data, padding...)
	_, err := d.OutEndpoint.Write(data)
	if err != nil {
		log.Fatalf("could not write data %d to device", data)
	}
}

// readFirmwareVersion reads the firmware version from `msg` and returns a formatted string
func (d *KrakenDriver) readFirmwareVersion(msg []byte) string {
	fwMajor, fwMinor, fwPatch := uint64(msg[0xb]), uint64(msg[0xc])<<8|uint64(msg[0xd]), uint64(msg[0xe])
	d.FirmwareVersion = []int{int(fwMajor), int(fwMinor), int(fwPatch)}
	d.CoolingProfiles = true

	return fmt.Sprintf("%d.%d.%d", fwMajor, fwMinor, fwPatch)
}

// setInstantSpeed sets a fixed speed per channel, but do not ensure persistence
func (d *KrakenDriver) setInstantSpeed(channel, duty string) {
	speedChannel, ok := speedChannels[channel]
	if !ok {
		log.Fatalf("channel %s not found", channel)
	}

	dutyInt, err := strconv.Atoi(duty)
	if err != nil {
		log.Fatal(err)
	}

	cbase, dmin, dmax := speedChannel[0], speedChannel[1], speedChannel[2]
	if dutyInt < dmin {
		dutyInt = dmin
	} else if dutyInt > dmax {
		dutyInt = dmax
	}

	d.write([]byte{0x2, 0x4d, byte(cbase & 0x70), 0, byte(dutyInt)})
}
