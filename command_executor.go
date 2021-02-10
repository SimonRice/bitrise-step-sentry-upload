package main

import "os/exec"

// CommandExecutor interface to allow mocking `exec.Command` within tests
type CommandExecutor interface {
	execute(string, ...string) ([]byte, error)
}

// StepExecutor implementation that Bitrise will use
type StepExecutor struct{}

/// Execute a console command
func (c StepExecutor) execute(command string, args ...string) ([]byte, error) {
	return exec.Command(command, args...).CombinedOutput()
}
