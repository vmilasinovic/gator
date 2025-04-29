package cli

import (
	"fmt"
	"reflect"
	"runtime"
)

func (c *Commands) Register(name, desc string, f func(*State, Command) error) {
	c.AllCommands[name] = f
	c.Descriptions[name] = desc

	// DEBUG
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	fmt.Printf("Registered handler: %v under name: %v", funcName, name)
}

func (c *Commands) RegisterCommands() {
	c.Register("login", "Log in to the specified user", handlerLogin)
}
