package main

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

var testConfig = Config{
	IsDebugMode:    "true",
	AuthToken:      "abcd12345",
	SentryURL:      "https://sentry.io/",
	OrgSlug:        "my-org",
	ProjectSlug:    "my-project",
	DsymPath:       "path/to/dsym",
	ProguardPath:   "path/to/proguard",
	ReleaseVersion: "myOrg/com.example.myApp@1.0.0",
}

/// TestCommandExecutor to mock command execution in tests
type TestCommandExecutor struct {
	ret []byte
	err error
}

func (c TestCommandExecutor) execute(command string, args ...string) ([]byte, error) {
	os.Args = args
	return c.ret, c.err
}

func TestDelegatePlatformUploads_Success(t *testing.T) {
	var tests = []struct {
		cmd      CommandExecutor
		cfg      Config
		expected []byte
	}{
		{
			cmd: TestCommandExecutor{
				ret: []byte("Success\n"),
				err: nil,
			},
			cfg: Config{
				SelectedPlatform: "both",
				IsDebugMode:      testConfig.IsDebugMode,
				AuthToken:        testConfig.AuthToken,
				SentryURL:        testConfig.SentryURL,
				OrgSlug:          testConfig.OrgSlug,
				ProjectSlug:      testConfig.ProguardPath,
				DsymPath:         testConfig.DsymPath,
				ProguardPath:     testConfig.ProguardPath,
			},
			expected: []byte("Uploads completed"),
		},
		{
			cmd: TestCommandExecutor{
				ret: []byte("Success\n"),
				err: nil,
			},
			cfg: Config{
				SelectedPlatform: "android",
				IsDebugMode:      testConfig.IsDebugMode,
				AuthToken:        testConfig.AuthToken,
				SentryURL:        testConfig.SentryURL,
				OrgSlug:          testConfig.OrgSlug,
				ProjectSlug:      testConfig.ProguardPath,
				ProguardPath:     testConfig.ProguardPath,
			},
			expected: []byte("Uploads completed"),
		},
		{
			cmd: TestCommandExecutor{
				ret: []byte("Success\n"),
				err: nil,
			},
			cfg: Config{
				SelectedPlatform: "ios",
				IsDebugMode:      testConfig.IsDebugMode,
				AuthToken:        testConfig.AuthToken,
				SentryURL:        testConfig.SentryURL,
				OrgSlug:          testConfig.OrgSlug,
				ProjectSlug:      testConfig.ProguardPath,
				DsymPath:         "mysd",
			},
			expected: []byte("Uploads completed"),
		},
	}

	for _, test := range tests {
		out, err := delegatePlatformUploads(test.cfg, test.cmd)
		if err != nil {
			t.Errorf("Test failed: %v", err)
		}
		if string(out) != string(test.expected) {
			t.Errorf("Test failed: expected %v but got %v", test.expected, out)
		}
		// reset Args
		os.Args = []string{}
	}
}

func TestDelegatePlatformUploads_Fail(t *testing.T) {
	var tests = []struct {
		cmd      CommandExecutor
		cfg      Config
		expected []byte
	}{
		{
			cmd: TestCommandExecutor{
				ret: []byte("Error\n"),
				err: errors.New("An error occurred"),
			},
			cfg: Config{
				SelectedPlatform: "linux",
			},
			expected: []byte{},
		},
		{
			cmd: TestCommandExecutor{
				ret: []byte("Error\n"),
				err: errors.New("An error occurred"),
			},
			cfg: Config{
				SelectedPlatform: "both",
			},
			expected: []byte("Error\n"),
		},
	}
	for _, test := range tests {
		out, err := delegatePlatformUploads(test.cfg, test.cmd)
		if err == nil {
			t.Errorf("%v: %v", err, string(out))
		}
		if string(out) != string(test.expected) {
			t.Errorf("Test failed: expected %v but got %v", test.expected, string(out))
		}
		// reset Args
		os.Args = []string{}
	}
}

func TestUploadSymbols_Success(t *testing.T) {
	var tests = []struct {
		cmd      CommandExecutor
		sentry   SentryCommand
		cfg      Config
		expected []string
	}{
		// proguard upload
		{
			cmd: TestCommandExecutor{
				ret: []byte("Success\n"),
				err: nil,
			},
			sentry: SentryCommand{
				Command:  uploadProguardCmd,
				FilePath: testConfig.ProguardPath,
			},
			cfg: testConfig,
			expected: []string{
				authTokenArg,
				testConfig.AuthToken,
				uploadProguardCmd,
				orgSlugArg,
				testConfig.OrgSlug,
				projectSlugArg,
				testConfig.ProjectSlug,
				testConfig.ProguardPath,
				logDebugArg,
			},
		},
		// dSYM upload
		{
			cmd: TestCommandExecutor{
				ret: []byte("Success\n"),
				err: nil,
			},
			sentry: SentryCommand{
				Command:  uploadDifCmd,
				FilePath: testConfig.DsymPath,
			},
			cfg: testConfig,
			expected: []string{
				authTokenArg,
				testConfig.AuthToken,
				uploadDifCmd,
				orgSlugArg,
				testConfig.OrgSlug,
				projectSlugArg,
				testConfig.ProjectSlug,
				testConfig.DsymPath,
				logDebugArg,
			},
		},
	}

	for _, test := range tests {
		_, err := uploadSymbols(test.cfg, test.sentry, test.cmd)
		if !reflect.DeepEqual(os.Args, test.expected) {
			t.Errorf("Test failed: Expected args %v, got %v", test.expected, os.Args)
		}
		if err != nil {
			t.Errorf("Test failed with error %v", err)
		}
		// reset Args
		os.Args = []string{}
	}
}

func TestUploadSymbols_Failed(t *testing.T) {
	// var cli =
	var tests = []struct {
		cmd      CommandExecutor
		sentry   SentryCommand
		cfg      Config
		expected error
	}{
		// proguard upload
		{
			cmd: TestCommandExecutor{
				ret: nil,
				err: errors.New("Upload failed"),
			},
			sentry: SentryCommand{
				Command:  uploadProguardCmd,
				FilePath: testConfig.ProguardPath,
			},
			cfg:      testConfig,
			expected: errors.New("Upload failed"),
		},
	}

	for _, test := range tests {
		_, err := uploadSymbols(test.cfg, test.sentry, test.cmd)
		if err == nil {
			t.Errorf("Test failed: Expected args %v, got %v", test.expected, os.Args)
		}
		// reset Args
		os.Args = []string{}
	}
}

func TestCreateFinalizeRelease_Success(t *testing.T) {
	const successString = "Success"

	cmd := TestCommandExecutor{ret: []byte(successString), err: nil}
	cfg := Config{
		IsDebugMode:    "true",
		AuthToken:      testConfig.AuthToken,
		OrgSlug:        testConfig.OrgSlug,
		ProjectSlug:    testConfig.ProjectSlug,
		ReleaseVersion: testConfig.ReleaseVersion,
	}
	// expected output from --auto `linkCommitsToRelease()`
	expectedArgs := []string{
		authTokenArg,
		testConfig.AuthToken,
		releasesCmd,
		orgSlugArg,
		testConfig.OrgSlug,
		setCommitsCmd,
		autoCommitsArg,
		testConfig.ReleaseVersion,
		logDebugArg,
	}
	expectedOut := successString

	out, err := createFinalizeRelease(cfg, cmd)
	if string(out) != expectedOut {
		t.Errorf("Test failed: Expected output %v, got %v", out, expectedOut)
	}
	if err != nil {
		t.Errorf("Test failed: Expected no error, got %v", err)
	}
	if !reflect.DeepEqual(os.Args, expectedArgs) {
		t.Errorf("Test failed: Expected args %v, got %v", expectedArgs, os.Args)
	} // reset Args
	os.Args = []string{}

}
func TestCreateFinalizeRelease_Failed(t *testing.T) {
	const missingVersionString = "Missing version string"
	const missingOrg = "Missing org slug"
	const missingProject = "Missing project slug"

	var tests = []struct {
		cmd          CommandExecutor
		cfg          Config
		expectedArgs []string
		expectedErr  error
	}{
		{
			cmd: TestCommandExecutor{ret: nil, err: nil},
			cfg: Config{
				IsDebugMode: "true",
				AuthToken:   testConfig.AuthToken,
				OrgSlug:     testConfig.OrgSlug,
				ProjectSlug: testConfig.ProjectSlug,
			},
			expectedArgs: os.Args,
			expectedErr:  nil,
		},
		{
			cmd: TestCommandExecutor{ret: nil, err: errors.New(missingOrg)},
			cfg: Config{
				IsDebugMode:    "true",
				AuthToken:      testConfig.AuthToken,
				OrgSlug:        "",
				ProjectSlug:    testConfig.ProjectSlug,
				ReleaseVersion: testConfig.ReleaseVersion,
			},
			expectedArgs: []string{
				authTokenArg,
				testConfig.AuthToken,
				releasesCmd,
				orgSlugArg,
				"", // missing org arg
				newReleaseSubCmd,
				projectSlugArg,
				testConfig.ProjectSlug,
				testConfig.ReleaseVersion,
				finalizeReleaseArg,
				logDebugArg,
			},
			expectedErr: errors.New(missingOrg),
		},
		{
			cmd: TestCommandExecutor{ret: nil, err: errors.New(missingProject)},
			cfg: Config{
				IsDebugMode:    "true",
				AuthToken:      testConfig.AuthToken,
				OrgSlug:        testConfig.OrgSlug,
				ProjectSlug:    "",
				ReleaseVersion: testConfig.ReleaseVersion,
			},
			expectedArgs: []string{
				authTokenArg,
				testConfig.AuthToken,
				releasesCmd,
				orgSlugArg,
				testConfig.OrgSlug,
				newReleaseSubCmd,
				projectSlugArg,
				"", // missing project arg
				testConfig.ReleaseVersion,
				finalizeReleaseArg,
				logDebugArg,
			},
			expectedErr: errors.New(missingProject),
		},
	}

	for _, test := range tests {
		_, err := createFinalizeRelease(test.cfg, test.cmd)
		if err == nil && test.expectedErr != nil {
			t.Errorf("Test failed: Expected error %v, got %v", test.expectedErr, err)
		}
		if !reflect.DeepEqual(os.Args, test.expectedArgs) {
			t.Errorf("Test failed: Expected args %v, got %v", test.expectedArgs, os.Args)
		}

		// reset Args
		os.Args = []string{}
	}
}

func TestLinkCommitsToRelease_Failed(t *testing.T) {
	const missingCommitsRef = "Repository does not exist"

	var tests = []struct {
		cmd          CommandExecutor
		cfg          Config
		expectedArgs []string
		expectedErr  error
	}{
		{
			cmd: TestCommandExecutor{ret: nil, err: errors.New(missingCommitsRef)},
			cfg: Config{
				IsDebugMode:       "true",
				AuthToken:         testConfig.AuthToken,
				OrgSlug:           testConfig.OrgSlug,
				ProjectSlug:       testConfig.ProjectSlug,
				AssociatedCommits: "myApp@123",
				ReleaseVersion:    testConfig.ReleaseVersion,
			},
			expectedArgs: []string{
				authTokenArg,
				testConfig.AuthToken,
				releasesCmd,
				orgSlugArg,
				testConfig.OrgSlug,
				setCommitsCmd,
				manualCommitsArg,
				"myApp@123",
				testConfig.ReleaseVersion,
				logDebugArg,
			},
			expectedErr: errors.New(missingCommitsRef),
		},
		{
			cmd: TestCommandExecutor{ret: nil, err: errors.New(missingCommitsRef)},
			cfg: Config{
				IsDebugMode: "true",
				AuthToken:   testConfig.AuthToken,
				OrgSlug:     testConfig.OrgSlug,
				ProjectSlug: testConfig.ProjectSlug,
			},
			expectedArgs: []string{
				authTokenArg,
				testConfig.AuthToken,
				releasesCmd,
				orgSlugArg,
				testConfig.OrgSlug,
				setCommitsCmd,
				autoCommitsArg,
				"",
				logDebugArg,
			},
			expectedErr: errors.New(missingCommitsRef),
		},
	}

	for _, test := range tests {
		_, err := linkCommitsToRelease(test.cfg, test.cmd)
		if err == nil {
			t.Errorf("Test failed: Expected error %v", test.expectedErr)
		}
		if !reflect.DeepEqual(os.Args, test.expectedArgs) {
			t.Errorf("Test failed: Expected args %v, got %v", test.expectedArgs, os.Args)
		} // reset Args
		os.Args = []string{}
	}
}
