package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"shipyard/internal/config"
	"shipyard/internal/privilege"
)

// distroName reads PRETTY_NAME from /etc/os-release, falls back to runtime.GOOS.
func distroName() string {
	f, err := os.Open("/etc/os-release")
	if err != nil {
		return runtime.GOOS
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "PRETTY_NAME=") {
			val := strings.TrimPrefix(line, "PRETTY_NAME=")
			val = strings.Trim(val, `"`)
			return val
		}
	}
	return runtime.GOOS
}

var checkCmd = &cobra.Command{
	Use:   "check",
	Short: "Check environment readiness (Docker, permissions, groups, etc.)",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("🔍 Shipyard Environment Check")
		fmt.Println("==============================")
		fmt.Println()

		issues := 0
		warnings := 0

		// ── System info ──
		u, _ := user.Current()
		fmt.Printf("👤 User:     %s\n", u.Username)
		fmt.Printf("💻 System:   %s/%s\n", distroName(), runtime.GOARCH)
		fmt.Println()

		// ── Docker daemon ──
		fmt.Println("── Docker ──")

		// Docker installed?
		fmt.Print("  Docker binary......... ")
		dockerPath, err := exec.LookPath("docker")
		if err != nil {
			fmt.Println("❌ not found in PATH")
			issues++
		} else {
			fmt.Printf("✅ %s\n", dockerPath)
		}

		// Docker daemon reachable?
		fmt.Print("  Docker daemon......... ")
		out, err := exec.Command("docker", "info", "--format", "{{.ServerVersion}}").Output()
		if err != nil {
			fmt.Println("❌ not reachable")
			issues++

			// Check if it's a permission issue
			if privilege.CanAccessDockerSocket() {
				fmt.Println("    → Socket is accessible but daemon may not be running")
				fmt.Println("    → Try: systemctl start docker")
			} else {
				fmt.Println("    → Cannot access Docker socket (/var/run/docker.sock)")
				fmt.Println("    → See Docker group check below")
			}
		} else {
			fmt.Printf("✅ v%s\n", strings.TrimSpace(string(out)))
		}

		// Docker without sudo?
		fmt.Print("  Docker rootless....... ")
		if privilege.IsRoot() {
			fmt.Println("⏭  running as root (N/A)")
		} else if privilege.CanAccessDockerSocket() {
			fmt.Println("✅ socket accessible without root")
		} else {
			fmt.Println("❌ socket not accessible without root")
			warnings++
			fmt.Println("    → Add yourself to the docker group:")
			fmt.Printf("    → sudo usermod -aG docker %s\n", u.Username)
			fmt.Println("    → Then log out and back in")
		}

		// Docker group membership
		fmt.Print("  Docker group.......... ")
		if privilege.IsUserInGroup("docker") {
			fmt.Println("✅ user is in 'docker' group")
		} else {
			fmt.Println("⚠  user is NOT in 'docker' group")
			warnings++
			fmt.Printf("    → sudo usermod -aG docker %s\n", u.Username)
		}

		// Docker Compose
		fmt.Print("  Docker Compose........ ")
		out, err = exec.Command("docker", "compose", "version", "--short").Output()
		if err != nil {
			fmt.Println("❌ not available")
			issues++
			fmt.Println("    → Install docker-compose or the compose plugin")
		} else {
			fmt.Printf("✅ v%s\n", strings.TrimSpace(string(out)))
		}
		fmt.Println()

		// ── Git ──
		fmt.Println("── Git ──")

		fmt.Print("  Git binary............ ")
		out, err = exec.Command("git", "version").Output()
		if err != nil {
			fmt.Println("❌ not found")
			issues++
		} else {
			fmt.Printf("✅ %s\n", strings.TrimSpace(string(out)))
		}

		// SSH agent
		fmt.Print("  SSH agent............. ")
		if os.Getenv("SSH_AUTH_SOCK") != "" {
			fmt.Println("✅ running")

			// Check if any keys are loaded
			out, err := exec.Command("ssh-add", "-l").Output()
			if err != nil {
				fmt.Println("    → No keys loaded (ssh-add -l returned error)")
				warnings++
			} else {
				lines := strings.Split(strings.TrimSpace(string(out)), "\n")
				fmt.Printf("    → %d key(s) loaded\n", len(lines))
			}
		} else {
			fmt.Println("⚠  not detected")
			warnings++
			fmt.Println("    → Private repo cloning may fail without SSH agent")
			fmt.Println("    → Try: eval $(ssh-agent) && ssh-add")
		}

		// SSH key exists
		fmt.Print("  SSH keys.............. ")
		home, _ := os.UserHomeDir()
		keyFiles := []string{"id_ed25519", "id_rsa", "id_ecdsa"}
		foundKeys := []string{}
		for _, kf := range keyFiles {
			if _, err := os.Stat(fmt.Sprintf("%s/.ssh/%s", home, kf)); err == nil {
				foundKeys = append(foundKeys, kf)
			}
		}
		if len(foundKeys) > 0 {
			fmt.Printf("✅ found: %s\n", strings.Join(foundKeys, ", "))
		} else {
			fmt.Println("⚠  no common SSH keys found in ~/.ssh/")
			warnings++
		}
		fmt.Println()

		// ── Yardctl Config ──
		fmt.Println("── Yardctl ──")

		fmt.Print("  Config file........... ")
		cfgPath := config.ConfigPath()
		if _, err := os.Stat(cfgPath); err != nil {
			fmt.Printf("❌ %s not found\n", cfgPath)
			fmt.Println("    → Run: yardctl init")
			issues++
		} else {
			fmt.Printf("✅ %s\n", cfgPath)
		}

		fmt.Print("  Data directory........ ")
		dataDir := config.DataDir()
		if _, err := os.Stat(dataDir); err != nil {
			fmt.Printf("❌ %s not found\n", dataDir)
			issues++
		} else {
			// Check if writable
			testFile := dataDir + "/.write_test"
			if err := os.WriteFile(testFile, []byte("test"), 0644); err != nil {
				fmt.Printf("⚠  %s exists but not writable\n", dataDir)
				warnings++
				fmt.Println("    → Run: yardctl init (to fix permissions)")
			} else {
				os.Remove(testFile)
				fmt.Printf("✅ %s (writable)\n", dataDir)
			}
		}

		fmt.Print("  Repos directory....... ")
		reposDir := config.ReposDir()
		if _, err := os.Stat(reposDir); err != nil {
			fmt.Printf("❌ %s not found\n", reposDir)
			issues++
		} else {
			fmt.Printf("✅ %s\n", reposDir)
		}

		// Systemd timer
		fmt.Print("  Sync timer............ ")
		out, err = exec.Command("systemctl", "is-active", "yardctl-sync.timer").Output()
		if err != nil {
			fmt.Println("⏹  not installed/active")
			fmt.Println("    → Install with: sudo systemctl enable --now yardctl-sync.timer")
		} else {
			status := strings.TrimSpace(string(out))
			if status == "active" {
				fmt.Println("✅ active")
			} else {
				fmt.Printf("⚠  %s\n", status)
			}
		}
		fmt.Println()

		// ── Summary ──
		fmt.Println("── Summary ──")
		if issues == 0 && warnings == 0 {
			fmt.Println("✅ All checks passed! Environment is ready.")
		} else if issues == 0 {
			fmt.Printf("⚠  %d warning(s), no critical issues. Environment should work.\n", warnings)
		} else {
			fmt.Printf("❌ %d issue(s), %d warning(s). Fix issues above before proceeding.\n", issues, warnings)
		}
	},
}

func init() {
	rootCmd.AddCommand(checkCmd)
}
