package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/vmilasinovic/gator.git/internal/cli"
	"github.com/vmilasinovic/gator.git/internal/config"
	"github.com/vmilasinovic/gator.git/internal/database"

	_ "github.com/lib/pq"
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

	// Open a connection to the database
	db, err := sql.Open("postgres", state.AppConfig.DBUrl)
	// Use generated database package to create a new *database.Queries
	dbQueries := database.New(db)
	state.Database = dbQueries

	// Create a context (for DB queries) that cancels a call if it takes longer than 3 seconds
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	state.Context = ctx

	// Start the CLI
	// cli.StartRepl(state, commands)

	// SHORT
	// Process command line args
	if len(os.Args) < 2 {
		fmt.Println("Error: Not enough arguments provided")
		os.Exit(1)
	}

	commandName := os.Args[1]
	args := os.Args[2:]

	// Check if command exists
	_, ok := commands.AllCommands[commandName]
	if !ok {
		fmt.Printf("Unknown command: %s\n", commandName)
		os.Exit(1)
	}

	// Create and run the command
	cmd := cli.Command{
		Name: commandName,
		Args: args,
	}

	if err := commands.Run(state, cmd); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
