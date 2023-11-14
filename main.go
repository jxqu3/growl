package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type GrowlCommand struct {
	Name        string
	Description string
	Command     string
	Shell       string
	ShellArgs   string
}
type GrowlYaml struct {
	Commands []GrowlCommand
}

func IndexFunc(s []GrowlCommand, f func(GrowlCommand) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

func runGoCmd(cmdname string, args []string, cfg GrowlYaml) error {
	if len(cfg.Commands) > 0 {
		if slices.Contains(cfg.Commands, GrowlCommand{Name: cmdname}) {
			runCommand(args, cfg)
			return nil
		}
	}
	fmt.Println("Executing: go", cmdname, ".")
	args = append([]string{cmdname, "."}, args...)
	cmd := exec.Command("go", args...)
	if err := cmd.Run(); err != nil {
		return err
	}
	return nil
}

func printErr(msg ...string) {
	color.New(color.FgRed).Println(msg)
	os.Exit(1)
}

func printList(commands []GrowlCommand) {
	color.New(color.FgGreen).Println("Commands:")
	for _, c := range commands {
		fmt.Println("-", c.Name)
		fmt.Println("   ", c.Description)
		fmt.Println("   ", c.Command)
	}
}

func runCommand(args []string, cfg GrowlYaml) {
	idx := IndexFunc(cfg.Commands, func(c GrowlCommand) bool { return c.Name == args[0] })
	if idx == -1 {
		printList(cfg.Commands)
		printErr("Command not found!")
	}

	cfgCmd := cfg.Commands[idx]

	for i := range args[:1] {
		cfgCmd.Command = strings.ReplaceAll(cfgCmd.Command, "%"+fmt.Sprint(i+1), args[0])
	}

	var shell string

	switch runtime.GOOS {
	case "windows":
		shell = "cmd"
	case "linux":
		shell = "bash"
	}

	if cfg.Commands[idx].Shell != "" {
		shell = cfg.Commands[idx].Shell
	}

	cmd := exec.Command(shell, cfg.Commands[idx].ShellArgs, cfg.Commands[idx].Command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		log.Fatal(err)
	}
}

func main() {
	var cfg GrowlYaml
	content, err := os.ReadFile("growl.yaml")
	if err != nil {
		log.Fatal(err)
	}
	yaml.Unmarshal(content, &cfg)
	app := cli.App{
		Name:                 "growl",
		Usage:                "simple go cli tools",
		EnableBashCompletion: true,
		Flags:                []cli.Flag{},
		Action: func(c *cli.Context) error {
			if len(c.Args().Slice()) == 0 {
				runGoCmd("run", c.Args().Slice(), cfg)
			} else {
				runCommand(c.Args().Slice(), cfg)
			}
			return nil
		},
		Commands: []*cli.Command{
			{
				Name: "run",
				Action: func(c *cli.Context) error {
					runGoCmd("run", c.Args().Slice(), cfg)
					return nil
				},
				Usage: "Run project (go run .)",
				Aliases: []string{
					"r",
				},
			},
			{
				Name: "build",
				Action: func(c *cli.Context) error {
					runGoCmd("build", c.Args().Slice(), cfg)
					return nil
				},
				Usage: "Build project (go build .)",
				Aliases: []string{
					"b",
				},
			},
		},
	}
	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
