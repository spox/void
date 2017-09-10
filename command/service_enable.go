package command

import (
	"fmt"
)

type ServiceEnableCommand struct {
	ServiceCommand
}

func (c *ServiceEnableCommand) Run(args []string) int {
	exitCode := 1
	if !c.isRoot() {
		c.UI.Error("This command must be run as `root`!")
		return exitCode
	}
	cOpts, err := c.Init(args, true)
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to setup service command: %s", err))
		return exitCode
	}
	if !c.ServiceExists() {
		c.UI.Error(fmt.Sprintf(
			"Service `%s` does not exist!", c.ServiceName))
		return exitCode
	}
	if c.ServiceIsEnabled() {
		c.UI.Error(fmt.Sprintf(
			"Service `%s` is already enabled!", c.ServiceName))
		return exitCode
	}
	if err = c.EnableService(); err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to enable service: %s", err))
		return exitCode
	} else {
		c.UI.Info(fmt.Sprintf(
			"Enabled service: %s", c.ServiceName))
	}
	if _, ok := cOpts["start"]; ok && !c.ServiceIsRunning() {
		c.UI.Warn(fmt.Sprintf(
			"Starting service `%s`...", c.ServiceName))
		if c.StartService() {
			c.UI.Error(fmt.Sprintf(
				"Failed to start service `%s`!", c.ServiceName))
			return exitCode
		} else {
			c.UI.Info(fmt.Sprintf(
				"Started service: %s", c.ServiceName))
		}
	}
	return 0
}
