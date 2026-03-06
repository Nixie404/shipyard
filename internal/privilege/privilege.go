package privilege

import (
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strings"
	"syscall"
)

// Escalator holds the detected privilege escalation command.
type Escalator struct {
	Command string // "sudo", "doas", or "" if root
}

// Detect finds the best available privilege escalation tool.
func Detect() *Escalator {
	// Already root, no escalation needed
	if os.Geteuid() == 0 {
		return &Escalator{Command: ""}
	}

	// Prefer doas if available
	if path, err := exec.LookPath("doas"); err == nil && path != "" {
		return &Escalator{Command: "doas"}
	}

	// Fall back to sudo
	if path, err := exec.LookPath("sudo"); err == nil && path != "" {
		return &Escalator{Command: "sudo"}
	}

	return nil
}

// IsRoot returns true if we're already running as root.
func IsRoot() bool {
	return os.Geteuid() == 0
}

// Run executes a command with privilege escalation if needed.
func (e *Escalator) Run(name string, args ...string) error {
	var cmd *exec.Cmd

	if e.Command == "" {
		// Already root
		cmd = exec.Command(name, args...)
	} else {
		fullArgs := append([]string{name}, args...)
		cmd = exec.Command(e.Command, fullArgs...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// RunOutput executes a command with privilege escalation and returns output.
func (e *Escalator) RunOutput(name string, args ...string) ([]byte, error) {
	var cmd *exec.Cmd

	if e.Command == "" {
		cmd = exec.Command(name, args...)
	} else {
		fullArgs := append([]string{name}, args...)
		cmd = exec.Command(e.Command, fullArgs...)
	}

	cmd.Stdin = os.Stdin
	cmd.Stderr = os.Stderr
	return cmd.Output()
}

// MkdirAll creates a directory, escalating privileges if needed.
func (e *Escalator) MkdirAll(path string, perm os.FileMode) error {
	// Try without escalation first
	err := os.MkdirAll(path, perm)
	if err == nil {
		return nil
	}

	// Check if it's a permission error
	if !os.IsPermission(err) {
		return err
	}

	if e == nil {
		return fmt.Errorf("permission denied and no privilege escalation tool found (install sudo or doas)")
	}

	fmt.Printf("⚡ Need elevated privileges to create %s\n", path)
	return e.Run("mkdir", "-p", path)
}

// WriteFile writes content to a file, escalating privileges if needed.
func (e *Escalator) WriteFile(path string, data []byte, perm os.FileMode) error {
	// Try without escalation first
	err := os.WriteFile(path, data, perm)
	if err == nil {
		return nil
	}

	if !os.IsPermission(err) {
		return err
	}

	if e == nil {
		return fmt.Errorf("permission denied and no privilege escalation tool found (install sudo or doas)")
	}

	fmt.Printf("⚡ Need elevated privileges to write %s\n", path)
	// Use tee with privilege escalation to write the file
	cmd := exec.Command(e.Command, "tee", path)
	cmd.Stdin = strings.NewReader(string(data))
	cmd.Stderr = os.Stderr
	// Suppress stdout from tee
	cmd.Stdout = nil
	return cmd.Run()
}

// Chown changes ownership of a path, escalating privileges if needed.
func (e *Escalator) Chown(path string, uid, gid int) error {
	err := os.Chown(path, uid, gid)
	if err == nil {
		return nil
	}

	if !os.IsPermission(err) {
		return err
	}

	if e == nil {
		return fmt.Errorf("permission denied and no privilege escalation tool found")
	}

	return e.Run("chown", "-R", fmt.Sprintf("%d:%d", uid, gid), path)
}

// Label returns a display string for the escalation method.
func (e *Escalator) Label() string {
	if e == nil {
		return "none (install sudo or doas)"
	}
	if e.Command == "" {
		return "running as root"
	}
	return e.Command
}

// IsUserInGroup checks if the current user is in the specified group.
func IsUserInGroup(groupName string) bool {
	u, err := user.Current()
	if err != nil {
		return false
	}

	gids, err := u.GroupIds()
	if err != nil {
		return false
	}

	group, err := user.LookupGroup(groupName)
	if err != nil {
		return false
	}

	for _, gid := range gids {
		if gid == group.Gid {
			return true
		}
	}
	return false
}

// CanAccessDockerSocket checks if we can access the Docker socket without root.
func CanAccessDockerSocket() bool {
	socketPath := "/var/run/docker.sock"
	// R_OK=4, W_OK=2
	err := syscall.Access(socketPath, 4|2)
	return err == nil
}
