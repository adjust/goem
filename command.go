package main

var commands = []*Command{
	cmdInit,
	cmdList,
	cmdBundle,
	cmdBuild,
	cmdTest,
	cmdHelp,
}

type Command struct {
	Run  func(args []string)
	Name string
}

func (c *Command) Runnable() bool {
	return c.Run != nil
}
