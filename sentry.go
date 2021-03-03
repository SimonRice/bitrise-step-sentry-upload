package main

const sentryCli = "sentry-cli"

/// `sentry-cli` command to upload dSYM file
const uploadDifCmd = "upload-dif"

/// `sentry-cli` command to upload proguard mapping
const uploadProguardCmd = "upload-proguard"

/// `sentry-cli` arg to enable debug logs
const logDebugArg = "--log-level=debug"

// SentryCommand allows the upload function to send to execute either
// `upload-proguard` or `upload-dif`
type SentryCommand struct {
	Command  string
	FilePath string
}

/// Builds the sentry-cli command string with the given args
func buildSentryArgs(cfg Config, command string) []string {
	return []string{
		"--auth-token",
		cfg.AuthToken,
		command,
		"--org",
		cfg.OrgSlug,
		"--project",
		cfg.ProjectSlug,
	}
}
