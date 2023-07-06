package main

import (
	"fmt"
	"os"

	"github.com/speedcrash100/go-yadro-testtask/pkg"
)

// Events
const (
	CLIENT_ARRIVED   = 1
	CLIENT_TAKE_SEAT = 2
	CLIENT_WAITING   = 3
	CLIENT_LEFT      = 4
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Usage: program <file>")
		return
	}

	file_path := args[1]

	o, err := os.Open(file_path)
	if err != nil {
		fmt.Println(err)
	}

	app := pkg.NewApp(o, os.Stdout)

	if err := app.Process(); err != nil {
		fmt.Println(err)
	}

}
