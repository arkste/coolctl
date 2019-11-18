// coolctl – A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72).
// Copyright (c) 2019 Arkadius Stefanski

// Package cmd contains all CLI commands
package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/arkste/coolctl/driver"
)

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "coolctl",
	Short: "A cross-platform driver for NZXT Kraken X (X42, X52, X62 or X72)",
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute(version string) {
	rootCmd.Version = version
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func init() {
	cobra.OnInitialize()
	rootCmd.PersistentFlags().IntVarP(&driver.Debug, "debug", "d", 1, "debug level")
}
