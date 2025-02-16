package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"devpathpro/pkg/config"
	"devpathpro/pkg/registry"
	"devpathpro/pkg/ui/cli"
	"devpathpro/pkg/ui/gui"
)

// Initialize logging
func init() {
	// Set up logging
	logFile, err := os.OpenFile("devpathpro.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(logFile)
}

func main() {
	// Parse command line flags
	cliMode := flag.Bool("cli", false, "Run in CLI mode instead of GUI")
	flag.Parse()

	// Check for administrator privileges
	if !registry.IsAdmin() {
		fmt.Println("Administrator privileges required")
		fmt.Println("Please restart the program with administrator privileges")
		os.Exit(1)
	}

	// Initialize configuration
	cfg := &config.Configuration{
		LogFile:  "devpathpro.log",
		Programs: config.GetDefaultPrograms(),
	}

	// Run in CLI or GUI mode based on flag
	if *cliMode {
		// CLI mode
		cli := cli.NewCLI(cfg)
		cli.Run()
	} else {
		// GUI mode (default)
		gui := gui.NewGUI()
		gui.Run()
	}
}
