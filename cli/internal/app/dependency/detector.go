package dependency

import (
	"fmt"
)

// Detector defines the interface for detecting installed tools
type Detector interface {
	// IsInstalled checks if a tool is installed
	// Returns: installed (bool), path (string), error
	IsInstalled(tool string) (bool, string, error)

	// GetInstallInstructions returns installation instructions for a tool
	GetInstallInstructions(tool string) string
}

// MultiDetector checks multiple sources for tool installation
type MultiDetector struct {
	detectors []Detector
}

// NewMultiDetector creates a new multi-source detector
func NewMultiDetector() *MultiDetector {
	return &MultiDetector{
		detectors: []Detector{
			&BrewDetector{},
			&NpmDetector{},
			&SystemDetector{},
		},
	}
}

// IsInstalled checks if a tool is installed using multiple detection methods
func (d *MultiDetector) IsInstalled(tool string) (bool, string, error) {
	// Special handling for markdownlint - try multiple variants
	if tool == "markdownlint" {
		// Try markdownlint-cli2 first (preferred, faster version)
		for _, detector := range d.detectors {
			installed, path, err := detector.IsInstalled("markdownlint-cli2")
			if err == nil && installed {
				return true, path, nil
			}
		}
		
		// Fall back to markdownlint-cli
		for _, detector := range d.detectors {
			installed, path, err := detector.IsInstalled("markdownlint-cli")
			if err == nil && installed {
				return true, path, nil
			}
		}
		
		// Try just "markdownlint" binary
		for _, detector := range d.detectors {
			installed, path, err := detector.IsInstalled("markdownlint")
			if err == nil && installed {
				return true, path, nil
			}
		}
		
		return false, "", nil
	}
	
	// Standard detection for other tools
	for _, detector := range d.detectors {
		installed, path, err := detector.IsInstalled(tool)
		if err != nil {
			continue
		}
		if installed {
			return true, path, nil
		}
	}
	return false, "", nil
}

// GetInstallInstructions returns formatted installation instructions
func (d *MultiDetector) GetInstallInstructions(tool string) string {
	instructions := fmt.Sprintf("Error: %s is not installed\n\n", tool)
	instructions += fmt.Sprintf("Marvin requires %s to run this check.\n\n", tool)
	instructions += "Installation options:\n\n"

	// Add tool-specific instructions
	switch tool {
	case "vale":
		instructions += "  Homebrew (recommended):\n"
		instructions += "    brew install vale\n\n"
		instructions += "  npm:\n"
		instructions += "    npm install -g vale\n\n"
		instructions += "  Manual:\n"
		instructions += "    https://vale.sh/docs/vale-cli/installation/\n"
	case "markdownlint", "markdownlint-cli", "markdownlint-cli2":
		instructions += "  npm (recommended):\n"
		instructions += "    npm install -g markdownlint-cli2\n\n"
		instructions += "  npm (alternative):\n"
		instructions += "    npm install -g markdownlint-cli\n\n"
		instructions += "  Homebrew:\n"
		instructions += "    brew install markdownlint-cli\n\n"
		instructions += "  Manual:\n"
		instructions += "    https://github.com/DavidAnson/markdownlint-cli2\n"
		instructions += "    https://github.com/igorshubovych/markdownlint-cli\n"
	default:
		instructions += fmt.Sprintf("  Please install %s and ensure it's in your PATH\n", tool)
	}

	instructions += "\nAfter installation, run this command again.\n"
	return instructions
}
