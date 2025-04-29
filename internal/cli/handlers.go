package cli

import (
	"encoding/json"
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 1 {
		s.AppConfig.CurrentUserName = cmd.Args[1]
		data, err := json.MarshalIndent(s.AppConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling updated config while logging in: %w", err)
		}

		if err = s.AppConfig.WriteToConf(s.AppConfig.FilePath, data); err != nil {
			return fmt.Errorf("error writing user to config while logging in: %w", err)
		}

		fmt.Printf("Logged in to user %v", s.AppConfig.CurrentUserName)
	} else {
		return fmt.Errorf("expected only 1 argument (username), please try again")
	}

	return nil
}
