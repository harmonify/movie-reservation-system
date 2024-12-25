package main

import (
	"os"

	"github.com/harmonify/movie-reservation-system/user-service/internal"
)

func main() {
	err := internal.StartApp()
	if err != nil {
		os.Exit(1)
	}
}
