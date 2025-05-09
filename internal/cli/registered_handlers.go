package cli

func (c *Commands) Register(name, desc string, f func(*State, Command) error) {
	c.AllCommands[name] = f
	c.Descriptions[name] = desc

	// DEBUG
	/*funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	trimmedName := path.Base(funcName)
	fmt.Printf("Registered handler: %v under name: %v\n", trimmedName, name)*/
}

func (c *Commands) RegisterCommands() {
	c.Register("login", "Log in to the specified user", handlerLogin)
	c.Register("register", "Register a new user", handlerRegister)
	c.Register("reset", "Clears the users table", handlerReset)
	c.Register("users", "List all usernames in the database", handlerGetUsers)
	c.Register("agg", "Fetch RSS feed from provided URL", handlerFetchRSS)
	c.Register("addfeed", "Add a new feed to DB", handlerAddFeed)
	c.Register("feeds", "List all fields from the DB", handlerFeeds)
	c.Register("follow", "Follow a feed with the provided URL", handlerFollow)
	c.Register("following", "List all feeds followed by user", handlerGetFeedFollowsForUser)
}
