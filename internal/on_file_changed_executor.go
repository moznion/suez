package internal

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/howeyc/fsnotify"
)

type OnFileChangedExecutor struct {
	CmdGenerator    func() *exec.Cmd
	WatchedFilePath string
	ExitOnError     bool
}

func (e *OnFileChangedExecutor) Execute() error {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return err
	}

	err = watcher.Watch(e.WatchedFilePath)
	if err != nil {
		return err
	}

	for {
		select {
		case <-watcher.Event:
			cmdOut, err := e.CmdGenerator().CombinedOutput()
			if err != nil {
				if e.ExitOnError {
					return fmt.Errorf("%w: %s", err, cmdOut)
				}
				_, _ = fmt.Fprintf(os.Stderr, "%s", cmdOut)
				continue
			}
			fmt.Printf("%s", cmdOut)
		case err := <-watcher.Error:
			return err
		}
	}
}
