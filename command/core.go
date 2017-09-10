package command

import (
	"fmt"
	"github.com/mitchellh/cli"
	"os"
	"os/exec"
	"strings"
	"syscall"
)

// Base command stubs
type CoreCommand struct {
	HelpText     string
	SynopsisText string
	Debug        bool
	UI           cli.Ui
}

func (c *CoreCommand) Help() string {
	return c.HelpText
}

func (c *CoreCommand) Synopsis() string {
	return c.SynopsisText
}

func (c *CoreCommand) Run(args []string) int {
	return 0
}

func Commands(ui cli.Ui, debug bool) map[string]cli.CommandFactory {
	return (&ServiceCommand{}).Commands(ui, debug)
}

// Extremely simple and dumb command line argument parser
func (c *CoreCommand) Parse(args []string) map[string][]string {
	fmtOpts := map[string][]string{}
	fmtOpts["_args_"] = []string{}
	for _, argItem := range args {
		if strings.Index(argItem, "-") == 0 {
			optVal := ""
			if strings.Contains(argItem, "=") {
				argParts := strings.SplitN(argItem, "=", 2)
				argItem = argParts[0]
				optVal = argParts[1]
			}
			optKey := strings.Replace(argItem, "-", "", 1)
			if _, ok := fmtOpts[optKey]; !ok {
				fmtOpts[optKey] = []string{}
			}
			if optVal != "" {
				fmtOpts[optKey] = append(fmtOpts[optKey], optVal)
			}
		} else {
			fmtOpts["_args_"] = append(fmtOpts["_args_"], argItem)
		}
	}
	return fmtOpts
}

// Runs the given command and returns the exit code. Includes
// debug information from the command execution.
func (c *CoreCommand) ExecuteCommand(cmd *exec.Cmd) int {
	exitCode := 1
	if c.Debug {
		if cmd.Stdout == nil {
			cmd.Stdout = os.Stdout
		}
		if cmd.Stderr == nil {
			cmd.Stderr = os.Stderr
		}
	}
	if err := cmd.Start(); err != nil {
		c.debug(fmt.Sprintf(
			"failed to start command `%x` - %s", cmd.Args, err))
		return exitCode
	}
	exitCode = 0
	if err := cmd.Wait(); err != nil {
		if exiterr, ok := err.(*exec.ExitError); ok {
			if status, ok := exiterr.Sys().(syscall.WaitStatus); ok {
				exitCode = status.ExitStatus()
				if exitCode != 0 {
					c.debug(fmt.Sprintf(
						"command returned non-zero exit: %d (%s)", exitCode, err))
				}
			}
		}
	}
	if exitCode == 0 && !cmd.ProcessState.Success() {
		c.debug("Exit code returned 0 but process state does " +
			"not show success. Setting exit code to 1.")
		exitCode = 1
	}
	return exitCode
}

func (c *CoreCommand) isRoot() bool {
	return os.Geteuid() == 0
}

func (c *CoreCommand) debug(line string) {
	if c.Debug {
		c.UI.Info(fmt.Sprintf(
			"[DEBUG] %s", line))
	}
}

func (c *CoreCommand) contains(haystack []string, needle string) bool {
	for _, v := range haystack {
		if v == needle {
			return true
		}
	}
	return false
}
