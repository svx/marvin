package dependency

import (
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// NpmDetector checks if tools are installed via npm
type NpmDetector struct{}

// IsInstalled checks if a tool is installed via npm (local or global)
func (d *NpmDetector) IsInstalled(tool string) (bool, string, error) {
	// Check local node_modules/.bin first
	localPath := filepath.Join("node_modules", ".bin", tool)
	if _, err := os.Stat(localPath); err == nil {
		absPath, _ := filepath.Abs(localPath)
		return true, absPath, nil
	}

	// Check if npm is available
	if _, err := exec.LookPath("npm"); err != nil {
		return false, "", err
	}

	// Check package.json for the tool
	if d.isInPackageJSON(tool) {
		return true, "package.json", nil
	}

	// Check global npm installation
	cmd := exec.Command("npm", "list", "-g", "--depth=0", "--json", tool)
	output, err := cmd.Output()
	if err != nil {
		return false, "", err
	}

	var result map[string]interface{}
	if err := json.Unmarshal(output, &result); err != nil {
		return false, "", err
	}

	if deps, ok := result["dependencies"].(map[string]interface{}); ok {
		if _, found := deps[tool]; found {
			// Get global npm bin path
			cmd = exec.Command("npm", "bin", "-g")
			binOutput, err := cmd.Output()
			if err != nil {
				return true, "", nil
			}
			binPath := strings.TrimSpace(string(binOutput))
			return true, filepath.Join(binPath, tool), nil
		}
	}

	return false, "", nil
}

// isInPackageJSON checks if a tool is listed in package.json
func (d *NpmDetector) isInPackageJSON(tool string) bool {
	data, err := os.ReadFile("package.json")
	if err != nil {
		return false
	}

	var pkg map[string]interface{}
	if err := json.Unmarshal(data, &pkg); err != nil {
		return false
	}

	// Check dependencies
	if deps, ok := pkg["dependencies"].(map[string]interface{}); ok {
		if _, found := deps[tool]; found {
			return true
		}
	}

	// Check devDependencies
	if devDeps, ok := pkg["devDependencies"].(map[string]interface{}); ok {
		if _, found := devDeps[tool]; found {
			return true
		}
	}

	return false
}

// GetInstallInstructions returns npm installation instructions
func (d *NpmDetector) GetInstallInstructions(tool string) string {
	return "npm install -g " + tool
}
