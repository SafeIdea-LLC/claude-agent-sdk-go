//go:build !windows

package transport

// platformArgs returns additional CLI arguments needed on this platform.
// On macOS/Linux, no extra args are needed — the CLI streams without --print.
func platformArgs() []string {
	return nil
}
