package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"

	"gh-release/internal/model"
)

type Config struct {
	Concurrency  int                `toml:"concurrency"`
	GitHub       ProviderConfig     `toml:"github"`
	GitLab       ProviderConfig     `toml:"gitlab"`
	Forgejo      ProviderConfig     `toml:"forgejo"`
	Repositories []RepositoryConfig `toml:"repositories"`
}

type ProviderConfig struct {
	Token string `toml:"token"`
}

type RepositoryConfig struct {
	Provider string `toml:"provider"`
	Host     string `toml:"host"`
	Owner    string `toml:"owner"`
	Name     string `toml:"name"`
}

func DefaultPath() string {
	home, err := os.UserHomeDir()
	if err != nil {
		return ".config/gh-release/config.toml"
	}

	return filepath.Join(home, ".config", "gh-release", "config.toml")
}

func Load(path string) (Config, error) {
	var cfg Config

	if _, err := toml.DecodeFile(path, &cfg); err != nil {
		return Config{}, fmt.Errorf("load config: %w", err)
	}

	if cfg.Concurrency < 1 {
		cfg.Concurrency = 4
	}

	return cfg, nil
}

func Save(path string, cfg Config) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config directory: %w", err)
	}

	f, err := os.OpenFile(
		path,
		os.O_WRONLY|os.O_CREATE|os.O_TRUNC,
		0600,
	)
	if err != nil {
		return fmt.Errorf("open config: %w", err)
	}
	defer f.Close()

	if err := toml.NewEncoder(f).Encode(cfg); err != nil {
		return fmt.Errorf("write config: %w", err)
	}

	if err := f.Chmod(0600); err != nil {
		return fmt.Errorf("set config permissions: %w", err)
	}

	return nil
}

func (c *Config) AddRepository(repo model.Repository) error {
	for _, existing := range c.Repositories {
		if existing.Provider == repo.Provider &&
			existing.Host == repo.Host &&
			existing.Owner == repo.Owner &&
			existing.Name == repo.Name {
			return fmt.Errorf("repository already exists: %s", repo.FullName())
		}
	}

	c.Repositories = append(c.Repositories, RepositoryConfig{
		Provider: repo.Provider,
		Host:     repo.Host,
		Owner:    repo.Owner,
		Name:     repo.Name,
	})

	return nil
}

func (c Config) Models() []model.Repository {
	repositories := make([]model.Repository, 0, len(c.Repositories))

	for _, repo := range c.Repositories {
		repositories = append(repositories, model.Repository{
			Provider: repo.Provider,
			Host:     repo.Host,
			Owner:    repo.Owner,
			Name:     repo.Name,
		})
	}

	return repositories
}
