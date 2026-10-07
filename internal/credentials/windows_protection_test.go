package credentials

import (
	"errors"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"runtime"
	"testing"
)

func TestWindowsProtectionFailsClosedForPrivateFileCredentials(t *testing.T) {
	if runtime.GOOS != "windows" {
		t.Skip("Windows native protection")
	}
	store := Store{Dir: t.TempDir()}
	for _, err := range []error{store.Put("private", "fixture"), func() error { _, e := store.Get("private"); return e }()} {
		var typed *output.Error
		if !errors.As(err, &typed) || typed.Code != 9 {
			t.Fatalf("unsupported protection must fail closed: %v", err)
		}
	}
}
