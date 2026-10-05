package cli

import (
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/controlplane"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
	"os"
	"strings"
)

func (a *App) requestCommands() {
	a.Root.AddCommand(&cobra.Command{Use: "request <method> <path>", Short: "Call a public API path with the selected credential", Args: cobra.ExactArgs(2), RunE: func(cmd *cobra.Command, args []string) error {
		method := strings.ToUpper(args[0])
		switch method {
		case "GET", "HEAD", "POST", "PUT", "PATCH", "DELETE":
		default:
			return output.New(2, "unsupported HTTP method")
		}
		b, e := a.body(false)
		if e != nil {
			return e
		}
		return a.call(cmd, method, args[1], b, false)
	}})
}
func (a *App) call(cmd *cobra.Command, method, path string, b []byte, secret bool) (err error) {
	defer func() {
		if err != nil && method != "GET" && method != "HEAD" && output.Normalize(err).Outcome == "" {
			err = notAttempted(err)
		}
	}()
	if e := controlplane.ValidatePath(path); e != nil {
		return e
	}
	c, e := a.client()
	if e != nil {
		return e
	}
	q, e := a.query()
	if e != nil {
		return e
	}
	if a.DryRun {
		return a.emit(map[string]any{"method": method, "path": path, "query": q, "body": json.RawMessage(b), "executed": false})
	}
	if method == "DELETE" && !a.Yes {
		return output.New(2, "destructive operation requires --yes")
	}
	var f *os.File
	if secret {
		if a.SecretFile == "" {
			return output.New(2, "issuance requires --secret-file")
		}
		f, e = os.OpenFile(a.SecretFile, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
		if e != nil {
			return output.New(2, "secret destination must not exist and must be writable")
		}
		defer f.Close()
	}
	v, _, e := c.Request(cmd.Context(), method, path, q, b)
	if e != nil {
		return e
	}
	if f != nil {
		enc := json.NewEncoder(f)
		if e = enc.Encode(v); e != nil {
			return &output.Error{Code: 10, Message: "issuance succeeded but saving secret failed", Outcome: "committed"}
		}
		if e = f.Sync(); e != nil {
			return &output.Error{Code: 10, Message: "issuance succeeded but secret sync failed", Outcome: "committed"}
		}
	}
	return a.emit(v)
}
