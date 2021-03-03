package main

import (
	"errors"
	"os"
	"reflect"
	"testing"
)

var testConfig = Config{
	IsDebugMode:  "true",
	AuthToken:    "abcd12345",
	SentryURL:    "https://sentry.io/",
	OrgSlug:      "my-org",
	ProjectSlug:  "my-project",
	DsymPath:     "path/to/dsym",
	ProguardPath: "path/to/proguard",
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
				"--auth-token",
				testConfig.AuthToken,
				uploadProguardCmd,
				"--org",
				testConfig.OrgSlug,
				"--project",
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
				"--auth-token",
				testConfig.AuthToken,
				uploadDifCmd,
				"--org",
				testConfig.OrgSlug,
				"--project",
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
