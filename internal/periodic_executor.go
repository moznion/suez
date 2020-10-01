package internal

import (
	"fmt"
	"os"
	"os/exec"
	"time"
)

type PeriodicExecutor struct {
	CmdGenerator func() *exec.Cmd
	Interval     time.Duration
	ExitOnError  bool
}

func (e *PeriodicExecutor) Execute() error {
	ticker := time.NewTicker(e.Interval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			cmdOut, err := e.CmdGenerator().CombinedOutput()
			if err != nil {
				if e.ExitOnError {
					return fmt.Errorf("%w: %s", err, cmdOut)
				}
				_, _ = fmt.Fprintf(os.Stderr, "%s", cmdOut)
				continue
			}
			fmt.Printf("%s", cmdOut)
		}
	}
}
