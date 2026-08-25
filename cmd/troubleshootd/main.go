package main

import (
	"fmt"
	"os"

	"callcentertroubleshooter/internal/app"
	"callcentertroubleshooter/internal/cli"
	"callcentertroubleshooter/internal/fixtures"
	"callcentertroubleshooter/internal/store"
	"callcentertroubleshooter/internal/troubleshoot"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, cli.RenderError(err))
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) == 0 {
		fmt.Println(cli.RenderHelp())
		return nil
	}
	command, err := cli.Parse(args)
	if err != nil {
		return fmt.Errorf("%w\n%s", err, cli.Usage())
	}
	if command.Name == "health" && command.Store == "" {
		return fmt.Errorf("store is required")
	}
	storage, err := store.Open(command.Store)
	if err != nil {
		return err
	}
	defer storage.Close()
	repository := fixtures.NewRepository()
	catalog := fixtures.NewCatalog(repository)
	queryer, err := troubleshoot.NewQueryer(catalog, troubleshoot.DefaultScopePolicy())
	if err != nil {
		return err
	}
	service, err := app.NewService(queryer, storage, catalog)
	if err != nil {
		return err
	}
	switch command.Name {
	case "diagnose":
		result, err := service.RunDiagnosisWorkflow(command.Employee, command.Actor)
		if err != nil {
			return err
		}
		fmt.Println(cli.RenderWorkflow(result))
	case "history":
		result, err := service.RunHistoryWorkflow(command.Employee)
		if err != nil {
			return err
		}
		fmt.Println(cli.RenderWorkflow(result))
	case "health":
		result, err := service.RunHealthWorkflow()
		if err != nil {
			return err
		}
		fmt.Println(cli.RenderWorkflow(result))
	}
	return nil
}
