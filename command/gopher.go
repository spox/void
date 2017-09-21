package command

import (
	"github.com/mitchellh/cli"
)

type GopherCommand struct {
	CoreCommand
}

func (c *GopherCommand) Commands(appName string, ui cli.Ui, debug bool) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"gopher convert": func() (cli.Command, error) {
			return &GopherConvertCommand{
				GopherCommand: GopherCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "Move project into GOPATH and link to current location",
						SynopsisText: "Move into GOPATH",
						UI:           ui,
						AppName:      appName,
						Flags: []CoreFlag{
							CoreFlag{
								Name:        "origin",
								Boolean:     false,
								Description: "Source origin",
								Default:     "github.com"}},
					},
				},
			}, nil
		},
	}
}
