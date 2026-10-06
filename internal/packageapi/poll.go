package packageapi

import (
	"context"
	"errors"
	"time"

	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
)

// Wait uses the caller's total deadline, including each HTTP request. Only GET
// observations are retried; acceptance, resume and cancellation are never replayed.
func (c *Client) Wait(ctx context.Context, initial Operation, observe func(Operation) error) (Operation, error) {
	current := initial
	for !current.Terminal {
		if err := ctx.Err(); err != nil {
			return current, output.Normalize(err)
		}
		pause := time.Duration(current.NextPollAfterMS) * time.Millisecond
		if pause < 100*time.Millisecond {
			pause = 100 * time.Millisecond
		}
		if pause > 30*time.Second {
			pause = 30 * time.Second
		}
		timer := time.NewTimer(pause)
		select {
		case <-ctx.Done():
			timer.Stop()
			return current, output.Normalize(ctx.Err())
		case <-timer.C:
		}
		next, err := c.Status(ctx, current.OperationID)
		if err != nil {
			if ctx.Err() != nil {
				return current, output.Normalize(ctx.Err())
			}
			var failure *output.Error
			if errors.As(err, &failure) && (failure.Code == 7 || failure.Code == 8) {
				continue
			}
			return current, err
		}
		if next.Revision < current.Revision || next.ArtifactDigest != current.ArtifactDigest {
			return current, output.New(9, "Package operation evidence regressed or changed identity")
		}
		current = next
		if observe != nil {
			if err := observe(current); err != nil {
				return current, err
			}
		}
	}
	return current, nil
}
