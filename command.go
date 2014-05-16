package main

var commands = []*Command{
	cmdInit,
	cmdList,
	cmdBundle,
	cmdBuild,
	cmdTest,
	cmdHelp,
	cmdMirrorConfig,
}

type Command struct {
	Run  func(args []string)
	Name string
}

func (c *Command) runnable() bool {
	return c.Run != nil
}
