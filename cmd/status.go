// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// Package cmd contains all CLI commands
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/arkste/coolctl/driver"
)

// statusCmd represents the status command
var statusCmd = &cobra.Command{
	Use:   "status",
	Short: "displays the current status",
	Run: func(cmd *cobra.Command, args []string) {
		kraken := driver.NewKrakenDriver()
		kraken.Connect()

		temperature, fanSpeed, pumpSpeed, firmwareVersion := kraken.GetStatus()

		fmt.Println(fmt.Sprintf("  Liquid temperature: %s °C", temperature))
		fmt.Println(fmt.Sprintf("  Fan speed: %d rpm", fanSpeed))
		fmt.Println(fmt.Sprintf("  Pump speed: %d rpm", pumpSpeed))
		fmt.Println(fmt.Sprintf("  Firmware Version: %s", firmwareVersion))
	},
}

func init() {
	rootCmd.AddCommand(statusCmd)
}
