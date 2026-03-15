package transport

import (
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
)

// resolveCLIForSubprocess checks whether the given CLI path is an npm wrapper
// script (e.g. node_modules/.bin/claude or claude.cmd) and, if so, resolves
// it to the underlying cli.js file invoked via "node". This is necessary
// because npm wrapper scripts fail when launched as subprocesses from Go on
// Windows: the .cmd wrapper has %dp0% path issues, and the shell wrapper is
// a bash script that can't execute under cmd.exe.
//
// Returns (executable, extraArgs) where:
//   - If resolved: executable="node", extraArgs=["path/to/cli.js"]
//   - If not resolved: executable=cliPath, extraArgs=nil
func resolveCLIForSubprocess(cliPath string) (string, []string) {
	// Normalize to forward slashes for consistent matching.
	norm := filepath.ToSlash(cliPath)

	// Only attempt resolution for paths that look like npm wrapper scripts.
	// These live in node_modules/.bin/ directories.
	if !isNpmWrapperPath(norm, cliPath) {
		return cliPath, nil
	}

	// Resolve cli.js relative to the wrapper.
	// npm wrappers in node_modules/.bin/ point to:
	//   ../\@anthropic-ai/claude-code/cli.js
	binDir := filepath.Dir(cliPath)
	cliJS := filepath.Join(binDir, "..", "@anthropic-ai", "claude-code", "cli.js")
	cliJS = filepath.Clean(cliJS)

	if _, err := os.Stat(cliJS); err != nil {
		// cli.js not found — fall back to original path.
		return cliPath, nil
	}

	// Find node executable.
	nodePath, err := exec.LookPath("node")
	if err != nil {
		// No node in PATH — can't use this approach, fall back.
		return cliPath, nil
	}

	return nodePath, []string{cliJS}
}

// isNpmWrapperPath returns true if the path looks like an npm wrapper script
// in a node_modules/.bin/ directory.
func isNpmWrapperPath(normalizedPath, originalPath string) bool {
	// Check if it's in a node_modules/.bin/ directory.
	if strings.Contains(normalizedPath, "node_modules/.bin/") {
		return true
	}

	// On Windows, also check for .cmd extension (npm generates these).
	if runtime.GOOS == "windows" {
		lower := strings.ToLower(originalPath)
		if strings.HasSuffix(lower, ".cmd") || strings.HasSuffix(lower, ".ps1") {
			return true
		}
	}

	// Check if the file is a shell script (starts with #!).
	// This catches cases where exec.LookPath finds a bash wrapper.
	if f, err := os.Open(originalPath); err == nil {
		defer f.Close()
		header := make([]byte, 2)
		if n, _ := f.Read(header); n == 2 && string(header) == "#!" {
			// It's a shell script — check if cli.js exists nearby.
			binDir := filepath.Dir(originalPath)
			cliJS := filepath.Join(binDir, "..", "@anthropic-ai", "claude-code", "cli.js")
			if _, err := os.Stat(filepath.Clean(cliJS)); err == nil {
				return true
			}
		}
	}

	return false
}
