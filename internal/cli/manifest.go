package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/manifest"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/spf13/cobra"
)

func (a *App) readManifest() (manifest.Document, error) {
	b, e := a.body(true)
	if e != nil {
		return manifest.Document{}, e
	}
	d, e := manifest.Parse(b)
	if e != nil {
		return d, output.New(2, e.Error())
	}
	return d, nil
}
func (a *App) validateManifest(d manifest.Document) error {
	for _, s := range d.Steps {
		var found *Operation
		for _, op := range a.Registry {
			if op.Command == s.Command {
				v := op
				found = &v
				break
			}
		}
		if found == nil {
			return output.New(2, "unknown manifest operation: "+s.Command)
		}
		if found.Status == "proposed" {
			return output.New(9, "manifest uses absent server extension")
		}
		if found.Secret || found.Effect == "publication" || found.Effect == "execution" || found.Method == "DELETE" {
			return output.New(2, "manifest permits configuration operations only; issuance publication execution and deletion must be explicit commands")
		}
		if len(s.Args) != len(found.Params) {
			return output.New(2, "wrong manifest resource arguments")
		}
		if found.Body && len(s.Body) == 0 {
			return output.New(2, "manifest operation requires body")
		}
	}
	return nil
}
func (a *App) manifestCommands() {
	g := a.group("manifest")
	g.AddCommand(&cobra.Command{Use: "validate", Args: cobra.NoArgs, RunE: func(*cobra.Command, []string) error {
		d, e := a.readManifest()
		if e != nil {
			return e
		}
		if e = a.validateManifest(d); e != nil {
			return e
		}
		return a.emit(map[string]any{"valid": true, "manifest_hash": d.Hash(), "authorization": "not_evaluated"})
	}})
	a.manifestPlanCommands(g)
}
