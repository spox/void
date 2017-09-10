package command

import (
	"fmt"
)

type ServiceDisableCommand struct {
	ServiceCommand
}

func (c *ServiceDisableCommand) Run(args []string) int {
	exitCode := 1
	if !c.isRoot() {
		c.UI.Error("This command must be run as `root`!")
		return exitCode
	}
	_, err := c.Init(args, true)
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
	if !c.ServiceIsEnabled() {
		c.UI.Error(fmt.Sprintf(
			"Service `%s` is not enabled!", c.ServiceName))
		return exitCode
	}
	if c.ServiceIsRunning() {
		c.UI.Warn(fmt.Sprintf(
			"Service `%s` is running. Stopping...", c.ServiceName))
		if !c.StopService() {
			c.UI.Error(fmt.Sprintf(
				"Failed to stop service `%s`!", c.ServiceName))
			return exitCode
		}
	}
	if err = c.DisableService(); err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to disable service: %s", err))
		return exitCode
	} else {
		c.UI.Info(fmt.Sprintf(
			"Disabled service: %s", c.ServiceName))
	}
	return 0
}
