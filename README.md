# coolctl

[![Go Report Card](https://goreportcard.com/badge/github.com/arkste/coolctl)](https://goreportcard.com/report/github.com/arkste/coolctl)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/arkste/coolctl/master/LICENSE)
[![Golang](https://img.shields.io/badge/Go-1.13-blue.svg)](https://golang.org)
![Linux](https://img.shields.io/badge/Supports-Linux-green.svg)

coolctl (cooler control) is a (soon) cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).

Currently still WIP.

## Installation

```bash
$ sudo apt-get install libusb-1.0 pkg-config
$ git clone https://github.com/arkste/coolctl.git
$ cd coolctl
$ make dep
```

## Get Status

```bash
$ go run main.go status
============================================
  Liquid temperature 32.7 °C
  Fan speed 527 rpm
  Pump speed 2040 rpm
  Firmware Version: 6.0.2
============================================
```

## Change Color

```bash
$ go run main.go color logo off
$ go run main.go color logo fading FF0000 00FF00 0000FF

$ go run main.go color ring off
$ go run main.go color ring fading FF0000 00FF00 0000FF
```

## Change Speed

```bash
$ go run main.go speed pump 20 60  35 60  55 100  60 100
$ go run main.go speed fan 20 25  35 25  50 55  60 100
```

## Full Silent Example

```bash
$ go run main.go color logo off
$ go run main.go color ring fading FF0000 00FF00 0000FF
$ go run main.go speed pump 20 60  35 60  55 100  60 100
$ go run main.go speed fan 20 25  35 25  50 55  60 100
```

## Credits

coolctl is a Golang-Port of [liquidctl](https://github.com/jonasmalacofilho/liquidctl) by [Jonas Malaco](https://github.com/jonasmalacofilho).