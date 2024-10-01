package Weather

import (
	"encoding/json"
	"fmt"
	"os"
)	

type Config struct {
	CitiesList      []string `json:"cities_list"`
	UpdateFrequency int      `json:"update_frequency"`
	CountDays       int      `json:"count_days"`
}

func (c* Config) Parse(path string) error {
	file, err := os.ReadFile(path)
	if err != nil {
		fmt.Printf("Error: File %s is not open!\n", path)
		os.Exit(1)
	}
	if err = json.Unmarshal(file, c); err != nil {
		fmt.Printf("Error: Config file is not valid!\n")
		os.Exit(1)
	}
	if len(c.CitiesList) == 0 {
		fmt.Printf("Error: Cities list is empty!\n")
		os.Exit(1)
	}
	if c.UpdateFrequency == 0 {
		c.UpdateFrequency = 1
  	}
	if c.CountDays == 0 {
		c.CountDays = 3
	}
	return nil
}
