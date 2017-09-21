package command

import (
	"errors"
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
	Flags        []CoreFlag
	UI           cli.Ui
	AppName      string
}

type CoreFlag struct {
	Description string
	Default     string
	Name        string
	Value       string
	Boolean     bool
}

type ParsedCli struct {
	Flags map[string]CoreFlag
	Args  []string
}

func (p *ParsedCli) Get(name string) *CoreFlag {
	flag, ok := p.Flags[name]
	if ok {
		return &flag
	}
	return nil
}

func (c *CoreCommand) Flag(name string) (CoreFlag, error) {
	var flag CoreFlag
	for _, flag := range c.Flags {
		if flag.Name == name {
			return flag, nil
		}
	}
	return flag, errors.New(fmt.Sprintf("Unknown flag `%s`", name))
}

func (c *CoreCommand) Help() string {
	maxLen := 0
	for _, flag := range c.Flags {
		if len(flag.Name) > maxLen {
			maxLen = len(flag.Name)
		}
	}
	maxLen++
	helpText := c.SynopsisText + "\n\nUsage: " + c.HelpText + "\n"
	for _, flag := range c.Flags {
		helpText = helpText + "    --" + flag.Name
		for i := 0; i < (maxLen - len(flag.Name)); i++ {
			helpText = helpText + " "
		}
		helpText = helpText + flag.Description + "\n"
	}
	return helpText
}

func (c *CoreCommand) Synopsis() string {
	return c.SynopsisText
}

func (c *CoreCommand) Run(args []string) int {
	return 0
}

func Commands(appName string, ui cli.Ui, debug bool) map[string]cli.CommandFactory {
	cmds := (&ServiceCommand{}).Commands(appName, ui, debug)
	for k, v := range (&GopherCommand{}).Commands(appName, ui, debug) {
		cmds[k] = v
	}
	return cmds
}

// Extremely simple and dumb command line argument parser
func (c *CoreCommand) Parse(args []string) (ParsedCli, error) {
	parsed := ParsedCli{
		Args:  []string{},
		Flags: map[string]CoreFlag{}}
	var lastItem CoreFlag
	setLast := false
	for _, argItem := range args {
		if !setLast {
			if strings.HasPrefix(argItem, "--") {
				argItem = strings.Replace(argItem, "--", "", 1)
			} else if strings.HasPrefix(argItem, "-") {
				argItem = strings.Replace(argItem, "-", "", 1)
			} else {
				parsed.Args = append(parsed.Args, argItem)
				continue
			}
			argParts := strings.SplitN(argItem, "=", 2)
			flag, err := c.Flag(string(argParts[0]))
			if err != nil {
				return parsed, err
			}
			if !flag.Boolean {
				if len(argParts) == 2 {
					flag.Value = string(argParts[1])
				} else {
					setLast = true
					lastItem = flag
				}
			}
			parsed.Flags[flag.Name] = flag
		} else {
			lastItem.Value = argItem
			setLast = false
		}
	}
	// Set any default flags that are unset
	for _, defFlag := range c.Flags {
		if _, ok := parsed.Flags[defFlag.Name]; !ok {
			if defFlag.Default != "" {
				defFlag.Value = defFlag.Default
				parsed.Flags[defFlag.Name] = defFlag
			}
		}
	}
	return parsed, nil
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
		c.UI.Output(fmt.Sprintf(
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
