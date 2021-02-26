package main

const sentryCli = "sentry-cli"

/// `sentry-cli` command to upload dSYM file
const uploadDifCmd = "upload-dif"

/// `sentry-cli` command to upload proguard mapping
const uploadProguardCmd = "upload-proguard"

/// `sentry-cli` command to create a new release
const releasesCmd = "releases"

/// `sentry-cli` `releases` subcommand
const newReleaseSubCmd = "new"

/// `sentry-cli` `reelases` subcommand to link commits to the given release
const setCommitsCmd = "set-commits"

/// `sentry-cli` arg to enable debug logs
const logDebugArg = "--log-level=debug"

/// `sentry-cli` required authorization token
const authTokenArg = "--auth-token"

/// `sentry-cli` organization slug argument
const orgSlugArg = "-o"

/// `sentry-cli` project slug argument
const projectSlugArg = "-p"

/// `sentry-cli` `release` command flag, finalizes the release on Sentry.io
const finalizeReleaseArg = "--finalize"

/// `sentry-cli` `release` command flag to automatically link commits
const autoCommitsArg = "--auto"

/// `sentry-cli` `release` command flag to manually link commits
const manualCommitsArg = "--commit"

// SentryCommand allows the upload function to send to execute either
// `upload-proguard` or `upload-dif`
type SentryCommand struct {
	Command  string
	FilePath string
}

/// Builds the sentry-cli command string with the given args
func buildSentryArgs(cfg Config, command string) []string {
	return []string{
		authTokenArg,
		cfg.AuthToken,
		command,
		orgSlugArg,
		cfg.OrgSlug,
		projectSlugArg,
		cfg.ProjectSlug,
	}
}

/// Builds the arg string for creating a new Sentry release
func buildReleaseArgs(cfg Config) []string {
	return []string{
		authTokenArg,
		cfg.AuthToken,
		releasesCmd,
		orgSlugArg,
		cfg.OrgSlug,
		newReleaseSubCmd,
		projectSlugArg,
		cfg.ProjectSlug,
		cfg.ReleaseVersion,
		finalizeReleaseArg,
	}
}

/// Builds the arg string for manually linking commits to the release
func linkManualCommitsArgs(cfg Config) []string {
	return []string{
		authTokenArg,
		cfg.AuthToken,
		releasesCmd,
		orgSlugArg,
		cfg.OrgSlug,
		setCommitsCmd,
		manualCommitsArg,
		cfg.AssociatedCommits,
		cfg.ReleaseVersion,
	}
}

/// Builds the arg string for auto linking commits to the release
func linkAutoCommitsArgs(cfg Config) []string {
	return []string{
		authTokenArg,
		cfg.AuthToken,
		releasesCmd,
		orgSlugArg,
		cfg.OrgSlug,
		setCommitsCmd,
		autoCommitsArg,
		cfg.ReleaseVersion,
	}
}
