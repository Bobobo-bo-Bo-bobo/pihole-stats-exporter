package main

import (
	"fmt"
	"runtime"
)

func showUsage() {
	showVersion()
	fmt.Printf(helpText, name)
}

func showVersion() {
	fmt.Printf(versionText, name, version, name, runtime.Version())
}
