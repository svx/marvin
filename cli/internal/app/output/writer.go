package output

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// Writer defines the interface for writing check results
type Writer interface {
	// Write saves results to the output directory
	Write(result *models.Result) (string, error)

	// GetOutputPath returns the full output path for a checker
	GetOutputPath(checkerName string) string
}

// JSONWriter writes results as JSON files
type JSONWriter struct {
	outputDir string
}

// NewJSONWriter creates a new JSON writer
func NewJSONWriter(outputDir string) *JSONWriter {
	if outputDir == "" {
		outputDir = ".marvin/results"
	}
	return &JSONWriter{
		outputDir: outputDir,
	}
}

// Write saves the result to a JSON file
func (w *JSONWriter) Write(result *models.Result) (string, error) {
	// Ensure output directory exists
	if err := os.MkdirAll(w.outputDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create output directory: %w", err)
	}

	// Generate filename with timestamp
	filename := w.GetOutputPath(result.Checker)

	// Marshal result to JSON with indentation
	data, err := json.MarshalIndent(result, "", "  ")
	if err != nil {
		return "", fmt.Errorf("failed to marshal result: %w", err)
	}

	// Write to file
	if err := os.WriteFile(filename, data, 0644); err != nil {
		return "", fmt.Errorf("failed to write file: %w", err)
	}

	return filename, nil
}

// GetOutputPath returns the full output path for a checker
func (w *JSONWriter) GetOutputPath(checkerName string) string {
	timestamp := time.Now().Format("20060102-150405")
	filename := fmt.Sprintf("%s-%s.json", checkerName, timestamp)
	return filepath.Join(w.outputDir, filename)
}
