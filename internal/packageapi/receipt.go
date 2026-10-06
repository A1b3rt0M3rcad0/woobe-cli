package packageapi

import (
	"bytes"
	"encoding/json"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/jsoninput"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/output"
	"github.com/A1b3rt0M3rcad0/woobe-cli/internal/packagecheckpoint"
)

type PlanReceipt struct {
	Format               string `json:"format"`
	SchemaVersion        string `json:"schema_version"`
	PrincipalFingerprint string `json:"principal_fingerprint"`
	Plan                 Plan   `json:"plan"`
}

func (r PlanReceipt) Validate() error {
	if r.Format != "woobe-package-plan-receipt" || r.SchemaVersion != "1.0" || !digestPattern.MatchString(r.PrincipalFingerprint) || !validExpiry(r.Plan.ExpiresAt) || !digestPattern.MatchString(r.Plan.PlanDigest) || !identifier.MatchString(r.Plan.PlanID) {
		return output.New(2, "Invalid or expired Package plan receipt")
	}
	return nil
}
func (r PlanReceipt) Save(path string) error {
	if err := r.Validate(); err != nil {
		return err
	}
	data, err := json.Marshal(r)
	if err != nil {
		return err
	}
	return packagecheckpoint.WritePrivateExclusive(path, data)
}
func LoadPlanReceipt(path string) (PlanReceipt, error) {
	var result PlanReceipt
	data, err := packagecheckpoint.ReadPrivate(path)
	if err != nil {
		return result, err
	}
	if jsoninput.Validate(data) != nil {
		return result, output.New(2, "Invalid Package plan receipt JSON")
	}
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.DisallowUnknownFields()
	if decoder.Decode(&result) != nil {
		return result, output.New(2, "Invalid Package plan receipt fields")
	}
	return result, result.Validate()
}
