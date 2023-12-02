package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/fatih/color"
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
	Extra       []string
}
type GrowlYaml struct {
	Shell     string
	GlobalEnv []GrowlEnv
	Commands  []GrowlCommand
}

func IndexFunc(s []GrowlCommand, f func(GrowlCommand) bool) int {
	for i := range s {
		if f(s[i]) {
			return i
		}
	}
	return -1
}

func printErr(msg ...string) {
	color.New(color.FgRed).Println("ERROR: ", strings.Join(msg, " "))
	os.Exit(1)
}

func isNumber(s string) bool {
	for _, c := range s {
		if c < '0' || c > '9' {
			return false
		}
	}
	return true
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
