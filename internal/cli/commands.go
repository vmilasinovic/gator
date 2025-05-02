package cli

import (
	"context"
	"fmt"

	"github.com/vmilasinovic/gator.git/internal/config"
	"github.com/vmilasinovic/gator.git/internal/database"
)

type State struct {
	AppConfig *config.Config
	Database  *database.Queries
	Context   context.Context
}

type Command struct {
	Name string
	Args []string
}

type Commands struct {
	AllCommands  map[string]func(*State, Command) error
	Descriptions map[string]string
}

func NewCommands() *Commands {
	appCommands := Commands{
		AllCommands:  make(map[string]func(*State, Command) error),
		Descriptions: make(map[string]string),
	}
	return &appCommands
}

func (c *Commands) Get() map[string]string {
	return c.Descriptions
}

func (c *Commands) Run(s *State, cmd Command) error {
	_, exists := c.AllCommands[cmd.Name]
	if exists {
		if err := c.AllCommands[cmd.Name](s, cmd); err != nil {
			return err
		}
	} else {
		return fmt.Errorf("invalid command")
	}
	return nil
}
