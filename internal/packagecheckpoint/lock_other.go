//go:build !linux && !darwin

package packagecheckpoint

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"os"
)

func supportedProtection() error {
	return output.New(9, "private package checkpoints require a supported OS protection mechanism")
}
func lockCheckpoint(path string) (*os.File, error) { return nil, supportedProtection() }

func openPrivateRead(path string) (*os.File, error) { return nil, supportedProtection() }
