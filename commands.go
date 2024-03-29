package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/pkg/errors"
	libcmd "github.com/spf13/cobra"
)

func init() {
	// flags
	rootCmd.PersistentFlags().BoolVarP(&debugFlag, "debug", "d", false, "output debug information")
	runCmd.PersistentFlags().BoolVarP(&quietFlag, "quiet", "q", false, "suppress stdout")
	runCmd.PersistentFlags().BoolVarP(&silentFlag, "silent", "s", false, "suppress both stdout and stderr")
	runCmd.PersistentFlags().BoolVarP(&lockdownFlag, "lockdown", "l", false, "disallow reading from stdin")
	// subcommands
	rootCmd.AddCommand(validateCmd)
	rootCmd.AddCommand(listDirectivesCmd)
	rootCmd.AddCommand(runCmd)
}

var debugFlag bool    // output all logs from the cmd app
var quietFlag bool    // suppress stdout but not stderr (useful for speed gains from output-heavy processes)
var silentFlag bool   // suppress both stdout and stderr (not recommended)
var lockdownFlag bool // does not read any input from stdin

var rootCmd = &libcmd.Command{
	Use:     command,
	Version: version,
}

var validateCmd = &libcmd.Command{
	Use:   "validate",
	Short: "Validate the file",
	Run: func(cmd *libcmd.Command, args []string) {

	},
}

var listDirectivesCmd = &libcmd.Command{
	Use:   "ls",
	Short: "List all directives",
	Run: func(cmd *libcmd.Command, args []string) {

	},
}

var runCmd = &libcmd.Command{
	Use:   "run",
	Short: "Run a command or script",
	Args:  libcmd.MinimumNArgs(1),
	PreRunE: func(cmd *libcmd.Command, args []string) error {
		return nil
	},
	RunE: func(cmd *libcmd.Command, args []string) error {
		cwd, _ := filepath.Abs(filepath.Dir(os.Args[0]))

		directivesFile, err := findDirectivesFile(cwd)
		if err != nil {
			return err
		}

		directives, err := parseDirectives(directivesFile)
		if err != nil {
			return err
		}

		exe, ok := directives[args[0]]
		if !ok {
			return errors.New("Unknown command")
		}
		fmt.Println(exe)

		return nil
	},
}

func findDirectivesFile(cwd string) (string, error) {
	ff, _ := ioutil.ReadDir(cwd)
	for _, f := range ff {
		if _, exists := validDirectivesFilenames[f.Name()]; exists && !f.IsDir() {
			return filepath.Abs(fmt.Sprintf("%s/%s", cwd, f.Name()))
		}
	}

	if cwd != "/" {
		parentDir, _ := filepath.Abs(fmt.Sprintf("%s/..", cwd))
		return findDirectivesFile(parentDir)
	}

	return "", errors.New("unable to find directives file")
}

func parseDirectives(fpath string) (map[string]directive, error) {
	var directives map[string]directive
	if _, err := toml.DecodeFile(fpath, &directives); err != nil {
		return nil, errors.Wrap(err, "failed to parse directives")
	}

	return directives, nil
}
