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
		// Can't get prefix, try to find the actual executable in PATH
		// For markdownlint-cli package, the executable is just "markdownlint"
		actualTool := tool
		if tool == "markdownlint-cli" || tool == "markdownlint-cli2" {
			actualTool = "markdownlint"
		}
		if path, err := exec.LookPath(actualTool); err == nil {
			return true, path, nil
		}
		return true, tool, nil
	}

	// The binary is typically in prefix/bin/tool
	prefix := strings.TrimSpace(string(output))
	
	// For markdownlint-cli package, the executable is just "markdownlint"
	binaryName := tool
	if tool == "markdownlint-cli" || tool == "markdownlint-cli2" {
		binaryName = "markdownlint"
	}
	
	path := prefix + "/bin/" + binaryName
	
	// Verify the binary exists
	if _, err := exec.LookPath(path); err != nil {
		// Try to find the actual executable in PATH
		if actualPath, err := exec.LookPath(binaryName); err == nil {
			return true, actualPath, nil
		}
		// Fall back to just the binary name (will be found in PATH)
		return true, binaryName, nil
	}
	
	return true, path, nil
}

// GetInstallInstructions returns Homebrew installation instructions
func (d *BrewDetector) GetInstallInstructions(tool string) string {
	return "brew install " + tool
}
