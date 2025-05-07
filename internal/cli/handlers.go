package cli

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/vmilasinovic/gator.git/internal/database"
	"github.com/vmilasinovic/gator.git/internal/rss"
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

func handlerFetchRSS(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		feed, err := rss.FetchFeed(s.Context, "https://www.wagslane.dev/index.xml")
		if err != nil {
			return fmt.Errorf("an error occured while fetching RSS feed: %w", err)
		}
		fmt.Println(feed)
	} else {
		return fmt.Errorf("agg command takes no arguments")
	}
	return nil
}

func handlerAddFeed(s *State, cmd Command) error {
	if len(cmd.Args) == 2 {
		name := cmd.Args[0]
		url := cmd.Args[1]

		currentUser, err := s.Database.GetUser(s.Context, s.AppConfig.CurrentUserName)
		if err != nil {
			log.Fatalf("an error occured while getting current user's details: %v", err)
		}
		currentUserID := currentUser.ID

		newFeed := database.AddFeedParams{
			UserID: currentUserID,
			Url:    url,
			Name:   name,
		}
		addedFeed, err := s.Database.AddFeed(s.Context, newFeed)
		if err != nil {
			return fmt.Errorf("an error occured while adding a new feed to DB: %w", err)
		}
		fmt.Println(addedFeed)

		addFollowing := Command{
			Name: "follow",
			Args: []string{},
		}
		addFollowing.Args = append(addFollowing.Args, url)
		handlerFollow(s, addFollowing)

	} else {
		return fmt.Errorf("addfeed takes 2 arguments - the name and the URL of the feed")
	}

	return nil
}

func handlerFeeds(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		feeds, err := s.Database.GetFeeds(s.Context)
		if err != nil {
			return fmt.Errorf("an error occured while fetching feeds from DB: %w", err)
		}

		for _, item := range feeds {
			fmt.Println(item)
		}
	} else {
		return fmt.Errorf("feeds take no arguments")
	}
	return nil
}

func handlerFollow(s *State, cmd Command) error {
	if len(cmd.Args) == 1 {
		url := cmd.Args[0]
		feedID, err := s.Database.GetFeedID(s.Context, url)
		if err != nil {
			if err == sql.ErrNoRows {
				return fmt.Errorf("no feed with that URL was found")
			}
			return fmt.Errorf("an error occured while checking the requested feed in the DB: %w", err)
		}

		checkUser, err := s.Database.GetUser(s.Context, s.AppConfig.CurrentUserName)
		if err != nil {
			return fmt.Errorf("an error occured while checking the requested user in the DB: %w", err)
		}
		user := checkUser.ID

		follow := database.InsertFeedFollowParams{
			UserID: user,
			FeedID: feedID,
		}
		newFollow, err := s.Database.InsertFeedFollow(s.Context, follow)
		if err != nil {
			return fmt.Errorf("an error occured while adding a new follow: %w", err)
		}
		fmt.Printf("%v is now following: %v\n", newFollow.UserName, newFollow.FeedName)

	} else {
		return fmt.Errorf("follow takes just 1 argument - URL")
	}

	return nil
}

func handlerGetFeedFollowsForUser(s *State, cmd Command) error {
	if len(cmd.Args) == 0 {
		checkUser, err := s.Database.GetUser(s.Context, s.AppConfig.CurrentUserName)
		if err != nil {
			return fmt.Errorf("an error occured while checking the requested user in the DB: %w", err)
		}
		user := checkUser.ID

		feedFollows, err := s.Database.GetFeedFollowsForUser(s.Context, user)
		if err != nil {
			return fmt.Errorf("an error occured while checking the requested feed in the DB: %w", err)
		}

		for _, feed := range feedFollows {
			fmt.Println(feed.FeedName)
		}

	} else {
		return fmt.Errorf("following takes no args")
	}

	return nil
}
