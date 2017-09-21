package command

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

type GopherConvertCommand struct {
	GopherCommand
}

func (c *GopherCommand) Run(args []string) int {
	exitCode := 1
	cmdOpts, err := c.Parse(args)
	cwd, err := os.Getwd()
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to determine current working directory: %s", err))
		return exitCode
	}
	projPath := cwd
	goPath := os.Getenv("GOPATH")
	if cmdOpts.Get("origin") == nil {
		c.UI.Error("Source origin is required for relocation!")
		return exitCode
	}
	originName := cmdOpts.Get("origin").Value
	relocateSelf := true
	if len(cmdOpts.Args) == 1 {
		projPath = cmdOpts.Args[0]
	} else if len(cmdOpts.Args) > 1 {
		c.UI.Error("Only single project can be converted at once.")
		return exitCode
	}
	projPath, err = filepath.Abs(projPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to expand project path: %s", err))
		return exitCode
	}
	relocateSelf = projPath == cwd
	// Check that path is not already link
	if _, err := os.Stat(projPath); err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to stat project directory: %s", err))
		return exitCode
	}

	if _, err := os.Readlink(projPath); err == nil {
		c.UI.Error(fmt.Sprintf(
			"Project directory is already a symlink: %s", projPath))
		return exitCode
	}

	// Check that path is a git repository
	_, err = os.Stat(filepath.Join(projPath, ".git"))
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Given project directory is not a git repository: %s", projPath))
		return exitCode
	}
	// Determine repository name, and container name
	pathList := strings.SplitN(projPath, string(filepath.Separator), -1)
	listLen := len(pathList)
	if listLen < 2 {
		c.UI.Error(fmt.Sprintf(
			"Cannot determine container directory from path: %s", projPath))
		return exitCode
	}
	repoName := pathList[listLen-1]
	ctnName := pathList[listLen-2]

	// Build our new path
	newProjPath := filepath.Join(goPath, "src", originName, ctnName, repoName)

	// Ensure destination does not already exist
	if _, err := os.Stat(newProjPath); err == nil {
		c.UI.Error(fmt.Sprintf(
			"Cannot relocate project. Destination already exists: %s", newProjPath))
		return exitCode
	}

	// Move if we need to relocate
	if relocateSelf {
		os.Chdir(filepath.Dir(projPath))
	}

	err = os.Rename(projPath, newProjPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to relocate project: %s", err))
		return exitCode
	}

	err = os.Symlink(newProjPath, projPath)
	if err != nil {
		c.UI.Error(fmt.Sprintf(
			"Failed to symlink project to original location: %s", err))
		return exitCode
	}

	if relocateSelf {
		os.Chdir(projPath)
	}
	exitCode = 0 // Successful conversion \o/
	c.UI.Info(fmt.Sprintf(
		"Successfully gophered project `%s/%s`!", ctnName, repoName))

	return exitCode
}
