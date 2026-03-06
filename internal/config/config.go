package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	DefaultConfigDir = "/etc/yardctl"
	DefaultDataDir   = "/var/lib/yardctl"
	ConfigFileName   = "config.json"
)

// RepoEntry represents a git repository tracked by yardctl.
type RepoEntry struct {
	URL    string `json:"url"`
	Branch string `json:"branch,omitempty"` // default: main
	Name   string `json:"name"`             // derived from URL
}

// StackEntry represents a discovered docker-compose stack.
type StackEntry struct {
	Name        string `json:"name"`
	RepoName    string `json:"repo_name"`
	ComposePath string `json:"compose_path"` // absolute path to docker-compose.yml
	Deployed    bool   `json:"deployed"`
}

// Config is the top-level yardctl configuration.
type Config struct {
	DataDir  string       `json:"data_dir"`
	ReposDir string       `json:"repos_dir"`
	Repos    []RepoEntry  `json:"repos"`
	Stacks   []StackEntry `json:"stacks"`
}

// ConfigDir returns the config directory.
func ConfigDir() string {
	dir := os.Getenv("YARDCTL_CONFIG_DIR")
	if dir == "" {
		dir = DefaultConfigDir
	}
	return dir
}

// ConfigPath returns the full path to the config file.
func ConfigPath() string {
	return filepath.Join(ConfigDir(), ConfigFileName)
}

// DataDir returns the data directory.
func DataDir() string {
	dir := os.Getenv("YARDCTL_DATA_DIR")
	if dir == "" {
		dir = DefaultDataDir
	}
	return dir
}

// ReposDir returns the repos directory inside data dir.
func ReposDir() string {
	return filepath.Join(DataDir(), "repos")
}

// Load reads the config from disk. If the file doesn't exist, returns a default config.
func Load() (*Config, error) {
	path := ConfigPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if os.IsNotExist(err) {
			return DefaultConfig(), nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}

	var cfg Config
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}
	return &cfg, nil
}

// Save writes the config to disk.
func (c *Config) Save() error {
	path := ConfigPath()
	dir := filepath.Dir(path)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create config dir: %w", err)
	}

	data, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config: %w", err)
	}
	return nil
}

// DefaultConfig returns a fresh default config.
func DefaultConfig() *Config {
	return &Config{
		DataDir:  DataDir(),
		ReposDir: ReposDir(),
		Repos:    []RepoEntry{},
		Stacks:   []StackEntry{},
	}
}

// FindRepo returns the repo entry with the given name, or nil.
func (c *Config) FindRepo(name string) *RepoEntry {
	for i := range c.Repos {
		if c.Repos[i].Name == name {
			return &c.Repos[i]
		}
	}
	return nil
}

// FindRepoByURL returns the repo entry with the given URL, or nil.
func (c *Config) FindRepoByURL(url string) *RepoEntry {
	for i := range c.Repos {
		if c.Repos[i].URL == url {
			return &c.Repos[i]
		}
	}
	return nil
}

// FindStack returns the stack entry with the given name, or nil.
func (c *Config) FindStack(name string) *StackEntry {
	// Try exact match first
	for i := range c.Stacks {
		if c.Stacks[i].Name == name {
			return &c.Stacks[i]
		}
	}

	// Try fuzzy suffix match if unique (e.g. "api" matching "repo-api")
	var candidates []*StackEntry
	for i := range c.Stacks {
		if strings.HasSuffix(c.Stacks[i].Name, "-"+name) {
			candidates = append(candidates, &c.Stacks[i])
		}
	}

	if len(candidates) == 1 {
		return candidates[0]
	}

	return nil
}

// RemoveRepo removes a repo by name and any stacks associated with it.
func (c *Config) RemoveRepo(name string) bool {
	found := false
	newRepos := []RepoEntry{}
	for _, r := range c.Repos {
		if r.Name == name {
			found = true
		} else {
			newRepos = append(newRepos, r)
		}
	}
	c.Repos = newRepos

	// Remove associated stacks
	newStacks := []StackEntry{}
	for _, s := range c.Stacks {
		if s.RepoName != name {
			newStacks = append(newStacks, s)
		}
	}
	c.Stacks = newStacks

	return found
}
