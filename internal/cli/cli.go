package cli

import (
	"fmt"
	"strings"
)

func Run(args []string) error {
	if len(args) == 0 {
		fmt.Println(renderHelp())
		return nil
	}
	command := strings.ToLower(strings.TrimSpace(args[0]))
	switch command {
	case "start":
		return runStart(args[1:])
	case "play":
		return runPlay(args[1:])
	case "history":
		return runHistory(args[1:])
	case "import":
		return runImport(args[1:])
	case "demo":
		return runDemo(args[1:])
	case "help", "-h", "--help":
		fmt.Println(renderHelp())
		return nil
	default:
		return fmt.Errorf("unknown command %q\n%s", args[0], renderHelp())
	}
}
