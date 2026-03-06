package stack

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"shipyard/internal/config"
)

// ComposeFileNames are the filenames we look for when discovering stacks.
var ComposeFileNames = []string{
	"docker-compose.yml",
	"docker-compose.yaml",
	"compose.yml",
	"compose.yaml",
}

// Discover walks through all repos and finds directories containing docker-compose files.
// Each directory with a compose file becomes a stack.
// Stacks are named as "repoName/relativePath" or just "repoName" if compose is at root.
func Discover(cfg *config.Config) ([]config.StackEntry, error) {
	var stacks []config.StackEntry

	for _, repo := range cfg.Repos {
		repoPath := filepath.Join(cfg.ReposDir, repo.Name)
		if _, err := os.Stat(repoPath); os.IsNotExist(err) {
			fmt.Printf("  ⚠ Repo directory missing for %s, skipping\n", repo.Name)
			continue
		}

		found, err := discoverInDir(repoPath, repo.Name)
		if err != nil {
			fmt.Printf("  ⚠ Error scanning %s: %v\n", repo.Name, err)
			continue
		}
		stacks = append(stacks, found...)
	}

	return stacks, nil
}

// discoverInDir walks a directory tree and finds compose files.
func discoverInDir(rootDir, repoName string) ([]config.StackEntry, error) {
	var stacks []config.StackEntry
	seen := make(map[string]bool)

	err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil // skip errors
		}

		// Skip hidden directories
		if info.IsDir() && strings.HasPrefix(info.Name(), ".") {
			return filepath.SkipDir
		}

		if info.IsDir() {
			return nil
		}

		for _, name := range ComposeFileNames {
			if info.Name() == name {
				dir := filepath.Dir(path)
				if seen[dir] {
					break
				}
				seen[dir] = true

				stackName := makeStackName(rootDir, dir, repoName)
				stacks = append(stacks, config.StackEntry{
					Name:        stackName,
					RepoName:    repoName,
					ComposePath: path,
					Deployed:    false,
				})
				break
			}
		}
		return nil
	})

	return stacks, err
}

// makeStackName creates a stack name from the repo name and the relative path.
func makeStackName(rootDir, composeDir, repoName string) string {
	rel, err := filepath.Rel(rootDir, composeDir)
	if err != nil || rel == "." {
		return repoName
	}
	// Replace path separators with dashes for a clean stack name
	rel = strings.ReplaceAll(rel, string(filepath.Separator), "-")
	return fmt.Sprintf("%s-%s", repoName, rel)
}

// Deploy brings up a stack using docker compose.
func Deploy(entry *config.StackEntry) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	fmt.Printf("🚀 Deploying stack %q from %s...\n", entry.Name, composeDir)

	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", entry.Name, "up", "-d")
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose up failed: %w", err)
	}

	entry.Deployed = true
	fmt.Printf("✅ Stack %q deployed successfully\n", entry.Name)
	return nil
}

// Destroy tears down a stack using docker compose.
func Destroy(entry *config.StackEntry) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	fmt.Printf("🛑 Destroying stack %q...\n", entry.Name)

	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", entry.Name, "down", "--remove-orphans")
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("docker compose down failed: %w", err)
	}

	entry.Deployed = false
	fmt.Printf("✅ Stack %q destroyed\n", entry.Name)
	return nil
}

// Logs shows logs for a stack or a specific service in a stack.
func Logs(entry *config.StackEntry, service string, follow bool) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	args := []string{"compose", "-f", composeFile, "-p", entry.Name, "logs"}
	if follow {
		args = append(args, "-f")
	}
	if service != "" {
		args = append(args, service)
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Ps lists containers in a stack.
func Ps(entry *config.StackEntry) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	cmd := exec.Command("docker", "compose", "-f", composeFile, "-p", entry.Name, "ps")
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Exec runs a command in a service container.
func Exec(entry *config.StackEntry, service string, command []string) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	args := []string{"compose", "-f", composeFile, "-p", entry.Name, "exec", service}
	args = append(args, command...)

	cmd := exec.Command("docker", args...)
	cmd.Dir = composeDir
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// Restart restarts services in a stack.
func Restart(entry *config.StackEntry, service string) error {
	composeDir := filepath.Dir(entry.ComposePath)
	composeFile := filepath.Base(entry.ComposePath)

	args := []string{"compose", "-f", composeFile, "-p", entry.Name, "restart"}
	if service != "" {
		args = append(args, service)
	}

	cmd := exec.Command("docker", args...)
	cmd.Dir = composeDir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// UpdateStacks merges newly discovered stacks into the config, preserving deploy state.
func UpdateStacks(cfg *config.Config, discovered []config.StackEntry) {
	// Build a map of existing stacks to preserve deployed state
	existing := make(map[string]*config.StackEntry)
	for i := range cfg.Stacks {
		existing[cfg.Stacks[i].Name] = &cfg.Stacks[i]
	}

	for i := range discovered {
		if old, ok := existing[discovered[i].Name]; ok {
			discovered[i].Deployed = old.Deployed
		}
	}

	cfg.Stacks = discovered
}
