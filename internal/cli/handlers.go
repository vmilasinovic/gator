package cli

import (
	"encoding/json"
	"fmt"
)

func handlerLogin(s *State, cmd Command) error {
	if len(cmd.Args) == 1 {
		requestedUser := cmd.Args[0]

		getUser, err := s.Database.GetUser(s.Context, requestedUser)
		if err != nil {
			return fmt.Errorf("an error occured while checking the requested user in the DB: %w", err)
		}

		s.AppConfig.CurrentUserName = getUser.Name

		data, err := json.MarshalIndent(s.AppConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling updated config while logging in: %w", err)
		}

		if err = s.AppConfig.WriteToConf(s.AppConfig.FilePath, data); err != nil {
			return fmt.Errorf("error writing user to config while logging in: %w", err)
		}

		fmt.Printf("Logged in to user %v\n", s.AppConfig.CurrentUserName)
	} else {
		return fmt.Errorf("expected only 1 argument (username), please try again")
	}

	return nil
}

func handlerRegister(s *State, cmd Command) error {
	if len(cmd.Args) == 1 {
		newUser := cmd.Args[0]

		createdUser, err := s.Database.CreateUser(s.Context, newUser)
		if err != nil {
			return fmt.Errorf("an error occured while registering a new user to the DB: %w", err)
		}

		fmt.Printf("Registered user %v\n", createdUser.Name)

		s.AppConfig.CurrentUserName = createdUser.Name
		data, err := json.MarshalIndent(s.AppConfig, "", "  ")
		if err != nil {
			return fmt.Errorf("error marshaling updated config while logging in: %w", err)
		}

		if err = s.AppConfig.WriteToConf(s.AppConfig.FilePath, data); err != nil {
			return fmt.Errorf("error writing user to config while logging in: %w", err)
		}

		fmt.Printf("Logged in to user %v\n", s.AppConfig.CurrentUserName)
	} else {
		return fmt.Errorf("expected only 1 argument (username), please try again")
	}

	return nil
}

func handlerReset(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		if err := s.Database.ClearUsers(s.Context); err != nil {
			return fmt.Errorf("an error occured while clearing table users: %w", err)
		}
	} else {
		return fmt.Errorf("reset command takes no arguments")
	}
	return nil
}

func handlerGetUsers(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		users, err := s.Database.GetUsers(s.Context)
		if err != nil {
			return fmt.Errorf("an error occured while fetching users from the DB: %w", err)
		}

		for _, user := range users {
			if user == s.AppConfig.CurrentUserName {
				fmt.Printf("* %v (current)\n", user)
			} else {
				fmt.Printf("* %v\n", user)
			}
		}
	} else {
		return fmt.Errorf("users command takes no arguments")
	}
	return nil
}
