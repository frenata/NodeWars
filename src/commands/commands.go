package commands

import (
	"argument"
	"fmt"
	"room"
	"strings"
)

type Command struct {
	Name string // strings that can be used to call this Command

	ShortDesc string
	LongDesc  string

	ArgsReq argument.ArgList
	ArgsOpt argument.ArgList

	Handler func(*client.Client, room.Room, []interface{}) error
}

func (c Command) Exec(cli *client.Client, context room.Room, strArgs []string) error {
	args, err := c.ValidateArgs(strArgs)
	if err != nil {
		// if we have trouble validating args
		return fmt.Errorf("%s\nusage: %s", err.Error(), c.Usage())
	} else {
		// otherwise actually execute the command
		fmt.Printf("<c.Exec> Calling command %s\n", c.Name)
		err = c.Handler(cli, context, args)
		if err != nil {
			return err
		}
	}
	return nil
}

// Help provides composed help info for the command
var padding int = 10

func (c Command) ShortHelp() string {
	// padding := strings.Repeat(" ", longest-len(c.Name))

	return fmt.Sprintf("%-*s - %s", padding, c.Name, c.ShortDesc)
}

func (c Command) LongHelp() string {
	ret := fmt.Sprintf("%s\nusage: %s", c.ShortHelp(), c.Usage())
	if c.LongDesc != "" {
		ret += "\n\n" + c.LongDesc
	}
	return ret
}

func (c Command) Usage() string {
	required := make([]string, len(c.ArgsReq))
	optional := make([]string, len(c.ArgsOpt))

	for i, arg := range c.ArgsReq {
		required[i] = arg.Name
	}

	for i, arg := range c.ArgsOpt {
		optional[i] = fmt.Sprintf("[%s]", arg.Name)
	}

	reqStr := strings.Join(required, " ")
	optStr := strings.Join(optional, " ")

	return fmt.Sprintf("%s %s %s", c.Name, reqStr, optStr)
}