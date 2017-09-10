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

func (c *ServiceCommand) Commands(ui cli.Ui, debug bool) map[string]cli.CommandFactory {
	return map[string]cli.CommandFactory{
		"service disable": func() (cli.Command, error) {
			return &ServiceDisableCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "Disable service and stop if running",
						SynopsisText: "Disable service",
						UI:           ui,
					},
				},
			}, nil
		},
		"service enable": func() (cli.Command, error) {
			return &ServiceEnableCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "Enable, and optionally start, a service\n  void service enable NAME [-start]",
						SynopsisText: "Enable service",
						UI:           ui,
					},
				},
			}, nil
		},
		"service list": func() (cli.Command, error) {
			return &ServiceListCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "List all services (or filtered)\n  void service list [-enabled] [-disabled]",
						SynopsisText: "List services",
						UI:           ui,
					},
				},
			}, nil
		},
		"service status": func() (cli.Command, error) {
			return &ServiceStatusCommand{
				ServiceCommand: ServiceCommand{
					CoreCommand: CoreCommand{
						Debug:        debug,
						HelpText:     "Display status of all or given services\n  void service status [NAME, NAME,...]",
						SynopsisText: "Status of services",
						UI:           ui,
					},
				},
			}, nil
		},
	}
}

func (c *ServiceCommand) Init(args []string, serviceName bool) (map[string][]string, error) {
	fmtOpts := c.Parse(args)
	if serviceName {
		if _, ok := fmtOpts["_args_"]; !ok || len(fmtOpts["_args_"]) != 1 {
			return fmtOpts, errors.New("Single service name required!")
		} else {
			c.ServiceName = fmtOpts["_args_"][0]
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
