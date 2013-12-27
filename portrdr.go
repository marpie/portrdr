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

	errCounter := 0

	// Setup TCP to TCP redirects
	for key, rdr := range config.Tcp2Tcp {
		fmt.Printf("[*] %s\n", key)
		if err := rdr.SetupRedirect(); err != nil {
			ErrorOut(err)
			errCounter += 1
		}
	}

	// Setup TCP to UDP redirects
	for key, rdr := range config.Tcp2Udp {
		fmt.Printf("[*] %s\n", key)
		if err := rdr.SetupRedirect(); err != nil {
			ErrorOut(err)
			errCounter += 1
		}
	}

	// Setup UDP to UDP redirects
	for key, rdr := range config.Udp2Udp {
		fmt.Printf("[*] %s\n", key)
		if err := rdr.SetupRedirect(); err != nil {
			ErrorOut(err)
			errCounter += 1
		}
	}

	// Setup UDP to TCP redirects
	for key, rdr := range config.Udp2Tcp {
		fmt.Printf("[*] %s\n", key)
		if err := rdr.SetupRedirect(); err != nil {
			ErrorOut(err)
			errCounter += 1
		}
	}

	// Check if any listeners are active
	if errCounter == config.Count() {
		ErrorExit(ERR_NO_REDIRECTIONS)
	}

	// Loop
	fmt.Println("\nPress Control+[C] to exit...")
	for {
		time.Sleep(time.Second)
	}
}
