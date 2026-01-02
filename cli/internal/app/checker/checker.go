package checker

import (
	"context"

	"github.com/svx/marvin/cli/internal/pkg/models"
)

// Checker defines the interface for all documentation checkers
type Checker interface {
	// Name returns the checker name
	Name() string

	// Check runs the check and returns results
	Check(ctx context.Context, opts CheckOptions) (*models.Result, error)

	// Validate validates the checker configuration
	Validate() error
}

// CheckOptions contains options for running a check
type CheckOptions struct {
	Path         string
	ConfigFile   string
	OutputFormat string
	ExtraArgs    []string
}
