package controlplane

// Only exact, fixed owner messages are classified. Never copy arbitrary server
// messages, input snapshots, Provider failures or protected values to output.
func packageExportIncompleteMessage(code, message string) string {
	if code != "PACKAGE_EXPORT_INCOMPLETE" {
		return ""
	}
	switch message {
	case "Snapshot has no complete ModelSpec":
		return "Package export blocked: selected snapshot has no complete ModelSpec; inspect the Agent Model configuration in the selected environment"
	case "Snapshot credential intent is unavailable":
		return "Package export blocked: selected snapshot has no Provider credential binding; inspect its Model and Provider configuration"
	case "Snapshot credential provider differs":
		return "Package export blocked: snapshot Model and credential use different Providers; inspect their bindings"
	case "Project Environment requirement is unavailable":
		return "Package export blocked: a required Project Environment field is unavailable; inspect the reported requirement path"
	}
	return ""
}
