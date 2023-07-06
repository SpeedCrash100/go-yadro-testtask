package main

import (
	"fmt"
	"os"
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

}
