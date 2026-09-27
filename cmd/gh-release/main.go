package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"

	"gh-release/internal/config"
	"gh-release/internal/provider"
	forgejoprovider "gh-release/internal/provider/forgejo"
	githubprovider "gh-release/internal/provider/github"
	gitlabprovider "gh-release/internal/provider/gitlab"
	"gh-release/internal/repository"
	"gh-release/internal/service"
	"gh-release/internal/ui"
)

func main() {
	if err := run(os.Args[1:]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}

func run(args []string) error {
	if len(args) > 0 {
		switch args[0] {
		case "add":
			return runAdd(args[1:])

		case "help", "--help", "-h":
			printUsage()
			return nil
		}
	}

	flags := flag.NewFlagSet("gh-release", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = printUsage

	configPath := flags.String(
		"config",
		config.DefaultPath(),
		"configuration file",
	)

	if err := flags.Parse(args); err != nil {
		return err
	}

	// Allow:
	//
	//   gh-release --config foo.toml add https://github.com/foo/bar
	//
	remaining := flags.Args()
	if len(remaining) > 0 {
		switch remaining[0] {
		case "add":
			return runAddWithConfig(*configPath, remaining[1:])
		case "help", "--help", "-h":
			printUsage()
			return nil
		default:
			return fmt.Errorf("unknown command %q", remaining[0])
		}
	}

	cfg, err := config.Load(*configPath)
	if err != nil {
		return err
	}

	registry := provider.NewRegistry()

	githubProvider, err := githubprovider.New(cfg.GitHub.Token)
	if err != nil {
		return fmt.Errorf("initialize github provider: %w", err)
	}

	registry.Register(githubProvider)
	registry.Register(gitlabprovider.New(cfg.GitLab.Token))
	registry.Register(forgejoprovider.New(cfg.Forgejo.Token))

	svc := service.New(registry, cfg.Concurrency)

	m := ui.New(svc, cfg.Models())

	if _, err := tea.NewProgram(m).Run(); err != nil {
		return err
	}

	return nil
}

func runAdd(args []string) error {
	return runAddWithConfig(config.DefaultPath(), args)
}

func runAddWithConfig(configPath string, args []string) error {
	flags := flag.NewFlagSet("gh-release add", flag.ContinueOnError)
	flags.SetOutput(os.Stdout)
	flags.Usage = addUsage

	if err := flags.Parse(args); err != nil {
		return err
	}

	args = flags.Args()

	if len(args) != 1 {
		flags.Usage()
		return fmt.Errorf("expected exactly one repository URL")
	}

	repo, err := repository.Parse(context.Background(), args[0])

	cfg, err := config.Load(configPath)
	if err != nil {
		return err
	}

	if err := cfg.AddRepository(repo); err != nil {
		return err
	}

	if err := config.Save(configPath, cfg); err != nil {
		return err
	}

	fmt.Printf(
		"added %s/%s (%s)\n",
		repo.Host,
		repo.FullName(),
		repo.Provider,
	)

	return nil
}

func printUsage() {
	fmt.Println(`gh-release - Git release tracker

Usage:
  gh-release [flags]
  gh-release add [flags] <repository-url>

Commands:
  add       Add a repository to track

Flags:
  -config string
            configuration file

Examples:
  gh-release
  gh-release add https://github.com/charmbracelet/bubbletea
  gh-release add https://gitlab.com/gitlab-org/gitlab
  gh-release add https://codeberg.org/forgejo/forgejo
  gh-release --config ~/.config/gh-release/config.toml`)
}

func addUsage() {
	fmt.Println(`Usage:
  gh-release add [flags] <repository-url>

Add a repository to the release tracker.

Supported repositories:
  https://github.com/owner/repository
  https://gitlab.com/owner/repository
  https://codeberg.org/owner/repository

Flags:`)
	fmt.Println("  -config string")
	fmt.Println("            configuration file")
}
