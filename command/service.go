package command

import (
	"errors"
	"fmt"
	"github.com/mitchellh/cli"
	"os"
	"os/exec"
	"path/filepath"
)

const SV_PATH = "/usr/bin/sv"
const SERVICES_PATH = "/etc/sv"
const ENABLED_SERVICES_PATH = "/var/service"

// Service command stub
type ServiceCommand struct {
	CoreCommand
	ServiceName string
}

func (c *ServiceCommand) Commands(appName string, ui cli.Ui, debug bool) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"service disable": func() (cli.Command, error) {
			return &ServiceDisableCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "void service disable NAME",
						SynopsisText: "Disable a system service",
						UI:           ui,
						AppName:      appName,
					},
				},
			}, nil
		},
		"service enable": func() (cli.Command, error) {
			return &ServiceEnableCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "void service enable NAME",
						SynopsisText: "Enable a system service",
						Flags: []CoreFlag{
							CoreFlag{
								Name:        "start",
								Boolean:     true,
								Description: "Start service after enabling"}},
						UI:      ui,
						AppName: appName,
					},
				},
			}, nil
		},
		"service list": func() (cli.Command, error) {
			return &ServiceListCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "void service list",
						SynopsisText: "List services",
						Flags: []CoreFlag{
							CoreFlag{
								Name:        "enabled",
								Boolean:     true,
								Description: "Display enabled services"},
							CoreFlag{
								Name:        "disabled",
								Boolean:     true,
								Description: "Display disabled services"}},

						UI:      ui,
						AppName: appName,
					},
				},
			}, nil
		},
		"service status": func() (cli.Command, error) {
			return &ServiceStatusCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "void service status [NAME,...]",
						SynopsisText: "Display status of services (all by default)",
						UI:           ui,
						AppName:      appName,
					},
				},
			}, nil
		},
	}
}

func (c *ServiceCommand) Init(args []string, serviceName bool) (ParsedCli, error) {
	fmtOpts, err := c.Parse(args)
	if err != nil {
		return fmtOpts, err
	}
	if serviceName {
		if len(fmtOpts.Args) != 1 {
			return fmtOpts, errors.New("Single service name required!")
		} else {
			c.ServiceName = fmtOpts.Args[0]
		}
	}
	return fmtOpts, nil
}

func (c *ServiceCommand) ServiceExists() bool {
	_, err := os.Stat(c.servicePath())
	return err == nil
}

func (c *ServiceCommand) ServiceIsEnabled() bool {
	_, err := os.Stat(c.enabledServicePath())
	return err == nil
}

func (c *ServiceCommand) ServiceIsRunning() bool {
	cmd := exec.Command(SV_PATH, "status", c.ServiceName)
	return c.ExecuteCommand(cmd) == 0
}

func (c *ServiceCommand) EnableService() error {
	return os.Symlink(c.servicePath(), c.enabledServicePath())
}

func (c *ServiceCommand) DisableService() error {
	return os.Remove(c.enabledServicePath())
}

func (c *ServiceCommand) StartService() bool {
	cmd := exec.Command(SV_PATH, "start", c.ServiceName)
	return c.ExecuteCommand(cmd) == 0
}

func (c *ServiceCommand) StopService() bool {
	cmd := exec.Command(SV_PATH, "stop", c.ServiceName)
	return c.ExecuteCommand(cmd) == 0
}

func (c *ServiceCommand) AllServices() ([]string, error) {
	return c.directoryList(SERVICES_PATH)
}

func (c *ServiceCommand) EnabledServices() ([]string, error) {
	return c.directoryList(ENABLED_SERVICES_PATH)
}

func (c *ServiceCommand) directoryList(path string) ([]string, error) {
	srvList := make([]string, 0)
	baseList, err := filepath.Glob(filepath.Join(path, "*"))
	if err != nil {
		return srvList, err
	}
	for _, v := range baseList {
		srvList = append(srvList, filepath.Base(v))
	}
	return srvList, nil
}

func (c *ServiceCommand) servicePath() string {
	return "/etc/sv/" + c.ServiceName
}

func (c *ServiceCommand) enabledServicePath() string {
	path, err := filepath.EvalSymlinks("/var/service")
	if err != nil {
		c.debug(fmt.Sprintf(
			"Failed to evaluate service path: %s", err))
		path = "/var/service"
	}
	return path + "/" + c.ServiceName
}
