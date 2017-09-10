package command

import (
	"fmt"
)

type ServiceStatusCommand struct {
	ServiceCommand
}

func (c *ServiceStatusCommand) Run(args []string) int {
	exitCode := 1
	if !c.isRoot() {
		c.UI.Error("This command must be run as `root`!")
		return exitCode
	}
	cOpts, _ := c.Init(args, true)
	eSrv, err := c.EnabledServices()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to list enabled services: %s", err))
		return exitCode
	}
	srvs := []string{}
	if len(cOpts["_args_"]) > 0 {
		srvs = cOpts["_args_"]
	} else {
		srvs = eSrv
	}
	for _, v := range srvs {
		if !c.contains(eSrv, v) {
			c.UI.Error(fmt.Sprintf(
				"Service name is not an enabled service: %s", v))
			return exitCode
		}
	}
	for _, v := range srvs {
		c.ServiceName = v
		if c.ServiceIsRunning() {
			c.UI.Info(v)
		} else {
			c.UI.Warn(v)
		}
	}
	return 0
}
