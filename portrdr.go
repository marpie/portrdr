package main

import (
	"fmt"
	"time"
)

func main() {
	config, err := LoadConfig(CONFIG_FILE)
	if err != nil {
		ErrorExit(err)
	}

	// Setup redirects
	errCounter := 0
	for _, rdr := range *config {
		if err := SetupRedirect(rdr); err != nil {
			ErrorOut(err)
			errCounter++
		}
	}
	// Check if any listeners are active
	if errCounter == len(*config) {
		ErrorExit(ERR_NO_REDIRECTIONS)
	}

	// Loop
	fmt.Println("\nPress Control+[C] to exit...")
	for {
		time.Sleep(time.Second)
	}
}
