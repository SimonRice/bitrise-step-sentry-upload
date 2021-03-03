package main

// Config will be populated with the retrieved values from environment variables
// configured as step inputs.
type Config struct {
	// Bitrise environment inputs
	SelectedPlatform string `env:"platform"`
	IsDebugMode      string `env:"is_debug_mode"`
	AuthToken        string `env:"auth_token"`
	SentryURL        string `env:"sentry_url"`
	OrgSlug          string `env:"org_slug"`
	ProjectSlug      string `env:"project_slug"`
	DsymPath         string `env:"dsym_path"`
	ProguardPath     string `env:"proguard_mapping_path"`
}
