package command

import (
	"errors"
	"github.com/Reece-Reklai/blog/internal/config/state"
)

type Command struct {
	Name      string
	Arguments []string
}

type CLI struct {
	Cmds map[string]func(*state.State, Command) error
}

func (cli *CLI) Run(states *state.State, cmd Command) error {
	value, ok := cli.Cmds[cmd.Name]
	if ok == false {
		return errors.New("Unknown Command")
	}
	err := value(states, cmd)
	return err
}

func (cli *CLI) Register(name string, handler func(states *state.State, cmd Command) error) error {
	_, ok := cli.Cmds[name]
	if ok == true {
		return errors.New("Command Exists")
	}
	cli.Cmds[name] = handler
	return nil
}
