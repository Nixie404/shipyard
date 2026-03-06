package repo

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"shipyard/internal/config"
)

// NameFromURL extracts a short repo name from a git URL.
// e.g. "git@github.com:org/deploy.git" -> "deploy"
// e.g. "https://github.com/org/deploy.git" -> "deploy"
func NameFromURL(url string) string {
	// Handle SSH-style URLs: git@github.com:org/repo.git
	if idx := strings.LastIndex(url, "/"); idx >= 0 {
		name := url[idx+1:]
		name = strings.TrimSuffix(name, ".git")
		return name
	}
	if idx := strings.LastIndex(url, ":"); idx >= 0 {
		name := url[idx+1:]
		if slashIdx := strings.LastIndex(name, "/"); slashIdx >= 0 {
			name = name[slashIdx+1:]
		}
		name = strings.TrimSuffix(name, ".git")
		return name
	}
	return strings.TrimSuffix(url, ".git")
}

// Add registers a new repository in the config and clones it.
func Add(cfg *config.Config, url string, branch string) (*config.RepoEntry, error) {
	if existing := cfg.FindRepoByURL(url); existing != nil {
		return nil, fmt.Errorf("repository %q already registered as %q", url, existing.Name)
	}

	name := NameFromURL(url)
	if existing := cfg.FindRepo(name); existing != nil {
		return nil, fmt.Errorf("a repository named %q already exists (url: %s)", name, existing.URL)
	}

	if branch == "" {
		branch = "main"
	}

	entry := config.RepoEntry{
		URL:    url,
		Branch: branch,
		Name:   name,
	}

	// Clone the repository
	repoPath := filepath.Join(cfg.ReposDir, name)
	if err := Clone(url, branch, repoPath); err != nil {
		return nil, fmt.Errorf("failed to clone: %w", err)
	}

	cfg.Repos = append(cfg.Repos, entry)
	return &entry, nil
}

// Clone performs a git clone.
func Clone(url, branch, dest string) error {
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return err
	}

	args := []string{"clone", "--branch", branch, "--single-branch", "--depth", "1", url, dest}
	cmd := exec.Command("git", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Pull performs a git pull on an existing repo.
func Pull(repoPath string) error {
	cmd := exec.Command("git", "-C", repoPath, "pull", "--ff-only")
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// SyncRepo pulls the latest for a single repo.
func SyncRepo(cfg *config.Config, entry *config.RepoEntry) error {
	repoPath := filepath.Join(cfg.ReposDir, entry.Name)

	// If dir doesn't exist, clone instead
	if _, err := os.Stat(repoPath); os.IsNotExist(err) {
		branch := entry.Branch
		if branch == "" {
			branch = "main"
		}
		fmt.Printf("  Cloning %s (%s)...\n", entry.Name, entry.URL)
		return Clone(entry.URL, branch, repoPath)
	}

	fmt.Printf("  Pulling %s...\n", entry.Name)
	return Pull(repoPath)
}

// SyncAll pulls all registered repos.
func SyncAll(cfg *config.Config) error {
	if len(cfg.Repos) == 0 {
		fmt.Println("No repositories registered. Use 'yardctl repo add <url>' first.")
		return nil
	}

	for i := range cfg.Repos {
		if err := SyncRepo(cfg, &cfg.Repos[i]); err != nil {
			fmt.Printf("  ⚠ Failed to sync %s: %v\n", cfg.Repos[i].Name, err)
		}
	}
	return nil
}

// Remove removes a repo directory from disk.
func RemoveDir(cfg *config.Config, name string) error {
	repoPath := filepath.Join(cfg.ReposDir, name)
	return os.RemoveAll(repoPath)
}
