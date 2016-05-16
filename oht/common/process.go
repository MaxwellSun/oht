package common

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type Process struct {
	Label     string
	ExecPath  string
	Args      []string
	Pid       int
	StartTime time.Time
	EndTime   time.Time
	Cmd       *exec.Cmd        `json:"-"`
	ExitState *os.ProcessState `json:"-"`
	WaitCh    chan struct{}    `json:"-"`
}

func StartProcess(label string, execPath string, args []string) (*Process, error) {
	cmd := exec.Command(execPath, args...)
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	proc := &Process{
		Label:     label,
		ExecPath:  execPath,
		Args:      args,
		Pid:       cmd.Process.Pid,
		StartTime: time.Now(),
		Cmd:       cmd,
		ExitState: nil,
		WaitCh:    make(chan struct{}),
	}
	go func() {
		err := proc.Cmd.Wait()
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				proc.ExitState = exitError.ProcessState
			}
		}
		proc.ExitState = proc.Cmd.ProcessState
		proc.EndTime = time.Now()
		if err != nil {
			fmt.Printf("Process: Error closing output file for %v: %v\n", proc.Label, err)
		}
		close(proc.WaitCh)
	}()
	return proc, nil
}

func (proc *Process) StopProcess(kill bool) error {
	if kill {
		return proc.Cmd.Process.Kill()
	} else {
		return proc.Cmd.Process.Signal(os.Interrupt)
	}
}
