package internal

import (
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/howeyc/fsnotify"
)

type PeriodicOnFileChangedExecutor struct {
	CmdGenerator    func() *exec.Cmd
	Interval        time.Duration
	WatchedFilePath string
	ExitOnError     bool
}

func (e *PeriodicOnFileChangedExecutor) Execute() error {
	ticker := time.NewTicker(e.Interval)
	defer ticker.Stop()

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Watch(e.WatchedFilePath)
	if err != nil {
		return err
	}

	fsEventChan := make(chan struct{}, 1)
	fsErrorChan := make(chan error, 0)
	go func() {
		for {
			select {
			case <-watcher.Event:
				select {
				case fsEventChan <- struct{}{}:
				default:
				}
			case err := <-watcher.Error:
				fsErrorChan <- err
			}
		}
	}()

	for {
		select {
		case <-ticker.C:
			cmdOut, ferr, err := func() ([]byte, error, error) {
				select {
				case err = <-fsErrorChan:
					return nil, err, nil
				case <-fsEventChan:
					cmdOut, err := e.CmdGenerator().CombinedOutput()
					if err != nil {
						return cmdOut, nil, fmt.Errorf("%w: %s", err, cmdOut)
					}
					return cmdOut, nil, nil
				default:
					return nil, nil, nil
				}
			}()

			if ferr != nil {
				return ferr
			}

			if err != nil {
				if e.ExitOnError {
					return err
				}
				_, _ = fmt.Fprintf(os.Stderr, "%s", cmdOut)
				continue
			}

			if cmdOut != nil {
				fmt.Printf("%s", cmdOut)
			}
		}
	}
}
