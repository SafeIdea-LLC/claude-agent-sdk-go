//go:build windows

package transport

import (
	"os/exec"
	"syscall"
)

// setProcAttr configures Windows-specific process creation flags.
// CREATE_NO_WINDOW (0x08000000) prevents the child process from inheriting
// or creating a console window. This matches the TS SDK's windowsHide: true
// and ensures Node.js detects !process.stdout.isTTY correctly, which is
// required for stream_event emission via the control protocol.
func setProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		CreationFlags: 0x08000000, // CREATE_NO_WINDOW
		HideWindow:    true,
	}
}
