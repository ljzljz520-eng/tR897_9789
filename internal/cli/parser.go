package cli

import (
	"fmt"
	"strings"
)

type Command struct {
	Name     string
	Employee string
	Actor    string
	Store    string
}

func Parse(args []string) (Command, error) {
	if len(args) == 0 {
		return Command{}, fmt.Errorf("command is required")
	}
	command := Command{Name: strings.ToLower(args[0]), Store: "troubleshooter.db"}
	for index := 1; index < len(args); index++ {
		switch args[index] {
		case "--employee":
			if index+1 >= len(args) {
				return Command{}, fmt.Errorf("--employee requires a value")
			}
			command.Employee = args[index+1]
			index++
		case "--actor":
			if index+1 >= len(args) {
				return Command{}, fmt.Errorf("--actor requires a value")
			}
			command.Actor = args[index+1]
			index++
		case "--store":
			if index+1 >= len(args) {
				return Command{}, fmt.Errorf("--store requires a value")
			}
			command.Store = args[index+1]
			index++
		default:
			return Command{}, fmt.Errorf("unknown argument %s", args[index])
		}
	}
	if command.Name != "diagnose" && command.Name != "history" && command.Name != "health" {
		return Command{}, fmt.Errorf("unknown command %s", command.Name)
	}
	return command, nil
}

func Usage() string {
	return "troubleshootd diagnose --employee 100001 --actor admin [--store path]\ntroubleshootd history [--employee 100001] [--store path]\ntroubleshootd health [--store path]"
}

func IsMutating(command Command) bool { return command.Name == "diagnose" }
