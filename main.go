package main

import (
	"log"

	"github.com/vmilasinovic/gator.git/internal/cli"
	"github.com/vmilasinovic/gator.git/internal/config"
)

func main() {
	// Initialize app configuration
	cfg, err := config.ReadGatorConfig()
	if err != nil {
		log.Fatalf("an error occured while reading ~/.gatorconfig.json: %v", err)
	}
	log.Printf("Config loaded successfully:\nDB URL: %v\nCurrent user: %v\n----------\n", cfg.DBUrl, cfg.CurrentUserName)
	state := &cli.State{AppConfig: cfg}

	// Initialize app commands
	commands := cli.NewCommands()
	commands.RegisterCommands()

	// Start the CLI
	cli.StartRepl(state, commands)
}
