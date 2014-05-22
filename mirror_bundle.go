package main

var cmdMirrorBundle = &Command{
	Run:  mirrorBundle,
	Name: "mirror-bundle",
}

func mirrorBundle(args []string) {
	makeBase()
	installDeps("", true)
	resolveDeps(args, true)
}
