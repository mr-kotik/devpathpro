package main

import (
	"fmt"
	"log"
	"os"

	"github.com/mr-kotik/devpathpro/pkg/config"
	"github.com/mr-kotik/devpathpro/pkg/registry"
	"github.com/mr-kotik/devpathpro/pkg/ui"
)

// Initialize logging
func init() {
	logFile, err := os.OpenFile("devpathpro.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Printf("Failed to open log file: %v", err)
		return
	}
	log.SetOutput(logFile)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}

func main() {
	// Check administrator privileges at startup
	if !registry.IsAdmin() {
		fmt.Println("\nPress Enter to exit...")
		fmt.Scanln()
		return
	}

	// Initialize configuration
	cfg := &config.Configuration{
		LogFile:  "devpathpro.log",
		Programs: config.GetDefaultPrograms(),
	}

	// Start the main menu
	ui.MainMenu(cfg)
} 