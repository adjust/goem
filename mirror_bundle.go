package main

var cmdMirrorBundle = &Command{
	Run:  mirrorBundle,
	Name: "mirror-bundle",
}

func mirrorBundle(args []string) {
	makeBase()
	installedPackages := installDeps("", true)
	resolveDeps(installedPackages, args, true)
}
