package controlplane

// Only exact, fixed owner messages are classified. Never copy arbitrary server
// messages, input snapshots, Provider failures or protected values to output.
func packageExportIncompleteMessage(code, message string) string {
	switch code {
	case "PACKAGE_EXPORT_NETWORK_BINDING":
		return "Network export blocked: an Agent or pinned snapshot is unavailable; inspect the reported node binding"
	case "PACKAGE_EXPORT_NETWORK_CONSTITUENT":
		return "Network export blocked: a constituent has no complete exact Agent snapshot; inspect the reported node and its selected Model/dependencies"
	case "PACKAGE_EXPORT_NETWORK_INTEGRITY":
		return "Network export blocked: saved snapshot integrity differs; preserve the request ID and investigate the original version"
	case "PACKAGE_EXPORT_NETWORK_DEFINITION":
		return "Network export blocked: saved definition is incompatible with the portable graph contract; inspect its topology and bindings"
	case "PACKAGE_EXPORT_NETWORK_CONFLICT":
		return "Network export blocked: saved snapshot or dependencies conflict; preserve the request ID for diagnosis"
	}
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
