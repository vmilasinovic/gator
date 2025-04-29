package cli

import (
	"fmt"
	"path"
	"reflect"
	"runtime"
)

func (c *Commands) Register(name, desc string, f func(*State, Command) error) {
	c.AllCommands[name] = f
	c.Descriptions[name] = desc

	// DEBUG
	funcName := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	trimmedName := path.Base(funcName)
	fmt.Printf("Registered handler: %v under name: %v\n", trimmedName, name)
}

func (c *Commands) RegisterCommands() {
	c.Register("login", "Log in to the specified user", handlerLogin)
}
