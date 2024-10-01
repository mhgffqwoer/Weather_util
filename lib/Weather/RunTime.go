package Weather

import (
	"os"
	"os/exec"
	"time"
	"golang.org/x/term"
)

func Update(args *Config, cityIdx *int, oldState *term.State) {
	for {
		time.Sleep(time.Duration(args.UpdateFrequency) * time.Minute)
		cmd := exec.Command("clear")
		cmd.Stdout = os.Stdout
		cmd.Run()
		weather := Weather{} 
		term.Restore(0, oldState)
		weather.Weather(args, cityIdx)
		_, _ = term.MakeRaw(0)
	}
}

func Run(args *Config) {
	weather := Weather{}
	cityIdx := 0

	cmd := exec.Command("clear")
	cmd.Stdout = os.Stdout
	cmd.Run()
	oldState, _ := term.GetState(0)
	weather.Weather(args, &cityIdx)
	_, _ = term.MakeRaw(0)

	go Update(args, &cityIdx, oldState)

	var key []byte = make([]byte, 1)
	for {
		os.Stdin.Read(key)
		if key[0] == 27 {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			term.Restore(0, oldState)
			os.Exit(0)
		}
		if key[0] == 'n' {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			cityIdx = (cityIdx + len(args.CitiesList) - 1) % len(args.CitiesList)
			term.Restore(0, oldState)
			weather.Weather(args, &cityIdx)
			_, _ = term.MakeRaw(0)
		}
		if key[0] == 'm' {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			cityIdx = (cityIdx + len(args.CitiesList) + 1) % len(args.CitiesList)
			term.Restore(0, oldState)
			weather.Weather(args, &cityIdx)
			_, _ = term.MakeRaw(0)
		}
		if key[0] == '-' {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			args.CountDays -= 1
			term.Restore(0, oldState)
			weather.Weather(args, &cityIdx)
			_, _ = term.MakeRaw(0)
		}
		if key[0] == '=' {
			cmd := exec.Command("clear")
			cmd.Stdout = os.Stdout
			cmd.Run()
			args.CountDays += 1
			term.Restore(0, oldState)
			weather.Weather(args, &cityIdx)
			_, _ = term.MakeRaw(0)
		}
	}
}
