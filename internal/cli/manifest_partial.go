package cli

import (
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)

type manifestPartial struct {
	Checkpoint             checkpoint
	Step, Message, Outcome string
	Saved                  bool
}

func (e *manifestPartial) Error() string { return e.Message }
func partialApply(cp checkpoint, step, message, outcome string, saved bool) error {
	return &manifestPartial{Checkpoint: cp, Step: step, Message: message, Outcome: outcome, Saved: saved}
}
func (a *App) emitManifestPartial(e *manifestPartial) int {
	counts := map[string]int{}
	for _, state := range e.Checkpoint.Steps {
		counts[state]++
	}
	data := map[string]any{"checkpoint": e.Checkpoint, "counts": counts, "stopped_at": e.Step, "checkpoint_saved": e.Saved}
	_ = output.Write(a.Out, a.Mode, output.Redact(data), map[string]string{"workspace_id": a.Workspace, "project_id": a.Project}, &output.Error{Code: 10, Message: e.Message, Outcome: e.Outcome})
	return 10
}
