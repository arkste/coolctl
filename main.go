package main

import "github.com/arkste/coolctl/driver"

func main() {
	kraken := driver.NewKrakenDriver()
	kraken.Connect()
	kraken.GetStatus()
	kraken.SetColor("ring", "loading", []string{"FF0000"})
	kraken.SetColor("logo", "pulse", []string{"FF0000"})
}
