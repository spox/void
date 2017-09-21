package main

import (
	"fmt"
	"github.com/chrisroberts/void/command"
	"github.com/mitchellh/cli"
	"os"
)

const VERSION = "0.1.0"
const APP_NAME = "void"

// Available commands
var Commands map[string]cli.CommandFactory

func main() {
	os.Exit(realMain())
}

func realMain() int {
	baseUi := &cli.BasicUi{Writer: os.Stdout, ErrorWriter: os.Stderr}
	ui := &cli.ColoredUi{
		ErrorColor:  cli.UiColorRed,
		InfoColor:   cli.UiColorGreen,
		OutputColor: cli.UiColorNone,
		WarnColor:   cli.UiColorYellow,
		Ui:          baseUi,
	}

	exitCode := 1
	debug := false
	if os.Getenv("VOID_DEBUG") != "" {
		debug = true
	}

	commands := command.Commands(APP_NAME, ui, debug)

	c := &cli.CLI{
		Args:     os.Args[1:],
		Commands: commands,
		Name:     APP_NAME,
		Version:  VERSION,
	}

	exitCode, err := c.Run()
	if err != nil {
		ui.Error(fmt.Sprintf("Unexpected error encountered: %s", err))
	}

	return exitCode
}
