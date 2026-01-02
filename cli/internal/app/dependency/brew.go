package dependency

import (
	"os/exec"
	"strings"
)

// BrewDetector checks if tools are installed via Homebrew
type BrewDetector struct{}

// IsInstalled checks if a tool is installed via Homebrew
func (d *BrewDetector) IsInstalled(tool string) (bool, string, error) {
	// Check if brew is available
	if _, err := exec.LookPath("brew"); err != nil {
		return false, "", err
	}

	// Check if the tool is installed via brew
	cmd := exec.Command("brew", "list", tool)
	if err := cmd.Run(); err != nil {
		return false, "", err
	}

	// Get the installation path
	cmd = exec.Command("brew", "--prefix", tool)
	output, err := cmd.Output()
	if err != nil {
		return true, tool, nil // Installed but can't get path, return tool name
	}

	// The binary is typically in prefix/bin/tool
	prefix := strings.TrimSpace(string(output))
	path := prefix + "/bin/" + tool
	
	// Verify the binary exists
	if _, err := exec.LookPath(path); err != nil {
		// Fall back to just the tool name (will be found in PATH)
		return true, tool, nil
	}
	
	return true, path, nil
}

// GetInstallInstructions returns Homebrew installation instructions
func (d *BrewDetector) GetInstallInstructions(tool string) string {
	return "brew install " + tool
}
