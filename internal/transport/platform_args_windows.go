//go:build windows

package transport

// platformArgs returns additional CLI arguments needed on Windows.
// On Windows, --print is required because the Ink TUI crashes on piped
// stdin ("Raw mode is not supported"), and without it the CLI falls back
// to a non-streaming mode that omits intermediate stream events.
func platformArgs() []string {
	return []string{"--print"}
}
