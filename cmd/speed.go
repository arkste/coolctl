package cmd

import (
	"errors"
	"fmt"
	"strings"

	"github.com/spf13/cobra"

	"github.com/arkste/coolctl/driver"
)

// speedCmd represents the speed command
var speedCmd = &cobra.Command{
	Use:   "speed",
	Short: "set the speed of the pump or fan",
	Args: func(cmd *cobra.Command, args []string) error {
		if len(args) < 1 {
			return errors.New("requires a speed channel (e.g: pump or fan)")
		}

		if len(args) < 2 {
			return errors.New("requires a speed profile (e.g: 20 25  35 25  50 55  60 100)")
		}

		return nil
	},
	Run: func(cmd *cobra.Command, args []string) {
		var profile string
		for i, profileNum := range args[1:] {
			profile += profileNum + " "
			if (i+1)%2 == 0 {
				profile += " "
			}
		}

		profile = strings.Trim(profile, " ")

		kraken := driver.NewKrakenDriver()
		kraken.Connect()
		kraken.SetSpeed(args[0], profile)

		fmt.Printf("%s: %s", args[0], profile)
		fmt.Println()
	},
}

func init() {
	rootCmd.AddCommand(speedCmd)
}
