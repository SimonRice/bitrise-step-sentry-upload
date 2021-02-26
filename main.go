package main

import (
	"errors"
	"fmt"
	"os"

	"github.com/bitrise-io/go-steputils/stepconf"
)

func delegatePlatformUploads(cfg Config, cmd CommandExecutor) ([]byte, error) {
	uploads := []SentryCommand{}
	dsym := SentryCommand{
		Command:  uploadDifCmd,
		FilePath: cfg.DsymPath,
	}
	proguard := SentryCommand{
		Command:  uploadProguardCmd,
		FilePath: cfg.ProguardPath,
	}
	if cfg.SelectedPlatform == "ios" {
		uploads = append(uploads, dsym)
	} else if cfg.SelectedPlatform == "android" {
		uploads = append(uploads, proguard)
	} else if cfg.SelectedPlatform == "both" {
		uploads = append(uploads, dsym, proguard)
	} else {
		return nil, errors.New("Error: selected_platform invalid")
	}

	for _, upload := range uploads {
		out, err := uploadSymbols(cfg, upload, cmd)
		if err != nil {
			return out, err
		}
		fmt.Printf("%s", out)
	}
	return []byte("Uploads completed"), nil
}

func uploadSymbols(cfg Config, sentry SentryCommand, cmd CommandExecutor) ([]byte, error) {
	args := buildSentryArgs(cfg, sentry.Command)
	args = append(args, sentry.FilePath)
	if cfg.IsDebugMode == "true" {
		args = append(args, logDebugArg)
	}

	fmt.Println(fmt.Sprintf("Executing %s, uploading %s...", sentry.Command, sentry.FilePath))
	return cmd.execute(sentryCli, args...)
}

func createFinalizeRelease(cfg Config, cmd CommandExecutor) ([]byte, error) {
	if cfg.ReleaseVersion == "" {
		return []byte("No release version declared, skipping Suspect Commit tracking..."), nil
	}
	args := buildReleaseArgs(cfg)
	if cfg.IsDebugMode == "true" {
		args = append(args, logDebugArg)
	}

	fmt.Println(fmt.Sprintf("Executing %s command, creating and finalizing release: %s", releasesCmd, cfg.ReleaseVersion))
	out, err := cmd.execute(sentryCli, args...)
	if err != nil {
		return out, err
	}
	fmt.Println(fmt.Sprintf("%s", out))
	return linkCommitsToRelease(cfg, cmd)
}

func linkCommitsToRelease(cfg Config, cmd CommandExecutor) ([]byte, error) {
	var args []string
	if cfg.AssociatedCommits != "" {
		args = append(args, linkManualCommitsArgs(cfg)...)
		fmt.Println(fmt.Sprintf("Manually linking %s, to release: %s...", cfg.AssociatedCommits, cfg.ReleaseVersion))
	} else {
		args = append(args, linkAutoCommitsArgs(cfg)...)
		fmt.Println(fmt.Sprintf("Automatically linking commits to release: %s...", cfg.ReleaseVersion))
	}
	if cfg.IsDebugMode == "true" {
		args = append(args, logDebugArg)
	}
	return cmd.execute(sentryCli, args...)
}

func main() {
	var cfg Config
	if err := stepconf.Parse(&cfg); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %s\n", err)
		os.Exit(1)
	}
	stepconf.Print(cfg)

	cmd := StepExecutor{}

	out, err := delegatePlatformUploads(cfg, cmd)
	if err != nil {
		fmt.Printf("%s\n", err)
		fmt.Printf("%s", string(out))
		os.Exit(1)
	}

	releaseOut, releaseErr := createFinalizeRelease(cfg, cmd)
	if releaseErr != nil {
		fmt.Printf("%s\n", releaseErr)
		fmt.Printf("%s", string(releaseOut))
		os.Exit(1)
	}
	fmt.Println(fmt.Sprintf("%s", releaseOut))

	os.Exit(0)
}
