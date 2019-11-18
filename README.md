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

## License

coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).  
Copyright (c) 2019 Arkadius Stefanski

coolctl is a Golang-Port of [liquidctl](https://github.com/jonasmalacofilho/liquidctl).  
Copyright (C) 2018–2019 [Jonas Malaco](https://github.com/jonasmalacofilho)  
Copyright (C) 2018–2019 each contribution's author  

This program is free software: you can redistribute it and/or modify it under the terms of the GNU General Public License as published by the Free Software Foundation, either version 3 of the License, or (at your option) any later version.

This program is distributed in the hope that it will be useful, but without any warranty; without even the implied warranty of MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the GNU General Public License for more details.

You should have received a copy of the GNU General Public License along with this program. If not, see https://www.gnu.org/licenses/.