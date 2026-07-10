package cli

import (
	"fmt"
)

func Run(args []string) error {
	if len(args) == 0 {
		printRootHelp()
		return nil
	}

	switch args[0] {
	case "-h", "--help", "help":
		printRootHelp()
		return nil
	case "-v", "--version", "version":
		fmt.Println(version)
		return nil
	case "config":
		return runConfig(args[1:])
	case "spec":
		return runSpec(args[1:])
	case "call":
		return runCall(args[1:])
	default:
		return fmt.Errorf("unknown command %q", args[0])
	}
}
