package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"slices"
	"strings"

	"github.com/fatih/color"
	"github.com/urfave/cli/v2"
	"gopkg.in/yaml.v2"
)

type GrowlEnv struct {
	Name  string
	Value string
}
type GrowlCommand struct {
	Name        string
	Description string
	Command     string
	Env         []GrowlEnv
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
		if slices.ContainsFunc(cfg.Commands, func(gc GrowlCommand) bool { return gc.Name == cmdname }) {
			args = append([]string{cmdname}, args...)
			runCommand(args, cfg)
			return nil
		}
	}
	fmt.Println("Executing: go", cmdname, ".", strings.Join(args, " "))
	args = append([]string{cmdname, "."}, args...)
	cmd := exec.Command("go", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

func printErr(msg ...string) {
	color.New(color.FgRed).Println("ERROR: ", strings.Join(msg, " "))
	os.Exit(1)
}

func printList(commands []GrowlCommand) {
	color.New(color.FgGreen).Println("Commands (growl.yaml):")
	blue := color.New(color.FgBlue)
	for _, c := range commands {
		blue.Print("- ")
		println(c.Name)
		blue.Print("    description ")
		fmt.Println(c.Description)
		blue.Print("    command ")
		fmt.Println(strings.Replace(c.Command, "%", blue.Sprint("%"), 1))
	}
}

func runCommand(args []string, cfg GrowlYaml) {
	idx := IndexFunc(cfg.Commands, func(c GrowlCommand) bool { return c.Name == args[0] })
	if idx == -1 {
		printList(cfg.Commands)
		printErr("Command not found!")
	}

	cfgCmd := cfg.Commands[idx]

	if len(args) > 1 {
		for i := range args[:1] {
			cfgCmd.Command = strings.ReplaceAll(cfgCmd.Command, "%"+fmt.Sprint(i+1), args[i+1])
		}
	}

	if cfgCmd.Shell == "" {
		switch runtime.GOOS {
		case "windows":
			cfgCmd.Shell = "cmd"
		case "linux", "darwin":
			cfgCmd.Shell = "sh"
		}
	}
	if cfgCmd.ShellArgs == "" {
		switch runtime.GOOS {
		case "windows":
			cfgCmd.ShellArgs = "/C"
		case "linux", "darwin":
			cfgCmd.ShellArgs = "-c"
		}
	}

	for _, v := range cfgCmd.Env {
		os.Setenv(v.Name, v.Value)
	}

	cmd := exec.Command(cfgCmd.Shell, cfgCmd.ShellArgs, cfgCmd.Command)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		printErr(err.Error())
	}
}

func initYaml() []byte {
	y, _ := yaml.Marshal(GrowlYaml{
		Commands: []GrowlCommand{
			{
				Name:        "hello",
				Description: "Says hello world!",
				Command:     "echo hello world, %1!",
			},
		},
	})
	os.WriteFile("growl.yaml", y, 0644)
	return y
}

func main() {
	var cfg GrowlYaml
	content, err := os.ReadFile("growl.yaml")
	if err != nil {
		fmt.Println("Growl.yaml not found. Generating it.")
		content = initYaml()
	}
	yaml.Unmarshal(content, &cfg)
	app := cli.App{
		Name:                 "growl",
		Usage:                "simple go cli tools",
		EnableBashCompletion: true,
		SkipFlagParsing:      true,
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
				SkipFlagParsing: true,
				Name:            "run",
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
				SkipFlagParsing: true,
				Name:            "init",
				Action: func(c *cli.Context) error {
					_, err = os.Stat("growl.yaml")
					if errors.Is(err, os.ErrNotExist) {
						initYaml()
						return nil
					}
					if err == nil {
						return errors.New("growl.yaml already exists")
					}
					return nil
				},
				Usage:   "Generates growl.yaml",
				Aliases: []string{},
			},
			{
				Name: "list",
				Action: func(c *cli.Context) error {
					printList(cfg.Commands)
					return nil
				},
				Usage: "List commands from growl.yaml",
				Aliases: []string{
					"l",
				},
			},
			{
				Name: "help",
				Action: func(c *cli.Context) error {
					if c.Args().Len() > 0 {
						for _, v := range c.App.Commands {
							if v.Name == c.Args().Get(0) {
								fmt.Println(v.UsageText)
							}
						}
						return nil
					}
					cli.ShowAppHelp(c)
					return nil
				},
				Usage: "Shows this help",
				Aliases: []string{
					"h",
				},
			},
			{
				Name:      "cross",
				UsageText: "growl cross --os [os] --arch [arch] --ldflags \"[ldflags]\" [--static] [--light] [--cgo] \ngrowl cross -o [os] -a [arch] -ld \"[ldflags]\" [-s] [-l] [-c]\nYou can use growl cross list to list available OS and CPU architectures",
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:    "os",
						Aliases: []string{"o"},
						Value:   runtime.GOOS,
					},
					&cli.StringFlag{
						Name:    "arch",
						Aliases: []string{"a"},
						Value:   runtime.GOARCH,
					},
					&cli.StringFlag{
						Name:    "ldflags",
						Aliases: []string{"ld"},
						Value:   "",
					},
					&cli.BoolFlag{
						Name:    "static",
						Aliases: []string{"s"},
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "light",
						Aliases: []string{"l"},
						Value:   false,
					},
					&cli.BoolFlag{
						Name:    "cgo",
						Aliases: []string{"c"},
						Value:   os.Getenv("CGO_ENABLED") == "1",
					},
				},
				Subcommands: []*cli.Command{
					{
						Name:  "list",
						Usage: "List available OS and CPU architectures",
						Action: func(c *cli.Context) error {
							color.Green("Available OS:")
							for _, v := range knownOS {
								fmt.Println(v)
							}
							color.Green("Available CPU architectures:")
							for _, v := range knownArch {
								fmt.Println(v)
							}
							return nil
						},
					},
				},
				Action: func(c *cli.Context) error {
					color.Green("Flags:")
					red := color.New(color.FgBlue)
					red.Print("- os ")
					os.Setenv("GOOS", c.String("os"))
					fmt.Println(c.String("os"))
					red.Print("- arch ")
					os.Setenv("GOARCH", c.String("arch"))
					fmt.Println(c.String("arch"))
					ld := c.String("ldflags")
					if c.Bool("static") {
						ld += " -extldflags=-static"
						os.Setenv("CGO_ENABLED", "1")
					}
					if c.Bool("light") {
						ld += " -w -s"
						os.Setenv("CGO_ENABLED", "1")
					}
					red.Print("- ldflags ")
					fmt.Println(ld)
					color.Green("Building...")
					args := append([]string{"build", "-ldflags=" + ld}, c.Args().Slice()...)
					cmd := exec.Command("go", args...)
					cmd.Stderr = os.Stderr
					cmd.Stdout = os.Stdout
					if err := cmd.Run(); err != nil {
						return err
					}
					return nil
				},
				Usage: "Build to target OS and arch. (growl cross --os=linux --arch=amd64)",
				Aliases: []string{
					"c",
				},
			},
			{
				SkipFlagParsing: true,
				Name:            "build",
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
		printErr(err.Error())
	}
}
