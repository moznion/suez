package suez

import (
	"errors"
	"os/exec"
	"time"

	"github.com/moznion/suez/internal"
)

type RunOpt struct {
	CmdGenerator    func() *exec.Cmd
	Interval        time.Duration
	WatchedFilePath string
	ExitOnError     bool
}

func (o *RunOpt) Validate() error {
	if o.Interval == 0 && o.WatchedFilePath == "" {
		return errors.New("either one of \"Interval\" or \"WatchedFilePath\" is a mandatory parameter, but both are missing")
	}
	return nil
}

func (o *RunOpt) getExecutor() (internal.Executor, error) {
	err := o.Validate()
	if err != nil {
		return nil, err
	}

	if o.Interval > 0 {
		if o.WatchedFilePath == "" {
			return &internal.PeriodicExecutor{
				CmdGenerator: o.CmdGenerator,
				Interval:     o.Interval,
				ExitOnError:  o.ExitOnError,
			}, nil
		}
		return &internal.PeriodicOnFileChangedExecutor{
			CmdGenerator:    o.CmdGenerator,
			Interval:        o.Interval,
			WatchedFilePath: o.WatchedFilePath,
			ExitOnError:     o.ExitOnError,
		}, nil
	}

	return &internal.OnFileChangedExecutor{
		CmdGenerator:    o.CmdGenerator,
		WatchedFilePath: o.WatchedFilePath,
		ExitOnError:     o.ExitOnError,
	}, nil
}

func Run(opt *RunOpt) error {
	executor, err := opt.getExecutor()
	if err != nil {
		return err
	}

	return executor.Execute()
}
