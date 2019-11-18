// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// coolctl is a Golang-Port of liquidctl.
// Copyright (C) 2018–2019 Jonas Malaco
// Copyright (C) 2018–2019 each contribution's author
package main

import "github.com/arkste/coolctl/cmd"

// Version represents the app version
var Version string

func main() {
	cmd.Execute(Version)
}
