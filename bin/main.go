package main

import (
	"Weather/lib/Weather"
	"flag"
	"fmt"
)

func main() {
	var ConfingPath string 
	flag.Usage = func() {
		fmt.Printf("Usage: %s [options]\n", "Weather")
		flag.PrintDefaults()
	}
	flag.StringVar(&ConfingPath, "f", "config.json", "get config file path [short]")
	flag.StringVar(&ConfingPath, "file", "config.json", "get config file path")
	flag.Parse()

	config := Weather.Config{}
	config.Parse(ConfingPath)

	Weather.Run(&config)


}
