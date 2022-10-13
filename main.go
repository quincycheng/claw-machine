package main

import (
	"fmt"
	"os"

	"github.com/quincycheng/claw-machine/app"
	"github.com/quincycheng/claw-machine/util"
)

func main() {
	config, err := util.LoadConfig()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)

	}
	app.Run(config)
}
