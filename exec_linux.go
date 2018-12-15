package main

import (
	"os/exec"
	"strings"
	"syscall"
)

// RunCommandOutput - output object for RunCommand
type RunCommandOutput struct {
	Command        string
	ReturnCode     int
	CombinedOutput []byte
}

// RunCommand - runs command
func RunCommand(name string, arg ...string) RunCommandOutput {

	var outputObj RunCommandOutput
	var err error

	// declare cmd object
	cmd := exec.Command(name, arg...)

	// set full command as string, for logging purposes
	if cmd.Args == nil || len(cmd.Args) == 0 {
		outputObj.Command = cmd.Path
	} else {
		outputObj.Command = strings.Join(cmd.Args, " ")
	}

	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid:   true,
		Pdeathsig: syscall.SIGKILL,
	}

	// execute command, acquire command return code
	outputObj.CombinedOutput, err = cmd.CombinedOutput()
	if err != nil {
		outputObj.ReturnCode = cmd.ProcessState.ExitCode()
	}

	return outputObj
}
