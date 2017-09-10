package command

import (
	"fmt"
)

type ServiceListCommand struct {
	ServiceCommand
}

func (c *ServiceListCommand) Run(args []string) int {
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
	if _, ok := cOpts["enabled"]; ok {
		for _, v := range eSrv {
			c.UI.Info(v)
		}
	} else {
		allSrv, err := c.AllServices()
		if err != nil {
			c.UI.Error(fmt.Sprintf(
				"Failed to list all services: %s", err))
			return exitCode
		}
		for _, v := range allSrv {
			if c.contains(eSrv, v) {
				if _, ok := cOpts["disabled"]; !ok {
					c.UI.Info(v)
				}
			} else {
				if _, ok := cOpts["enabled"]; !ok {
					c.UI.Error(v)
				}
			}
		}
	}
	return 0
}
