//go:build !windows

package transport

import "os/exec"

// setProcAttr is a no-op on non-Windows platforms.
func setProcAttr(cmd *exec.Cmd) {}
