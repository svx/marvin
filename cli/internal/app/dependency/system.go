package dependency

import (
	"os/exec"
)

// SystemDetector checks if tools are available in the system PATH
type SystemDetector struct{}

// IsInstalled checks if a tool is available in the system PATH
func (d *SystemDetector) IsInstalled(tool string) (bool, string, error) {
	path, err := exec.LookPath(tool)
	if err != nil {
		return false, "", err
	}
	return true, path, nil
}

// GetInstallInstructions returns generic installation instructions
func (d *SystemDetector) GetInstallInstructions(tool string) string {
	return "Please install " + tool + " and ensure it's in your PATH"
}
