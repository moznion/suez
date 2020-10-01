package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"time"

	"github.com/moznion/suez"
)

var (
	ver string
	rev string
)

func main() {
	var (
		intervalSec     int64
		watchedFilePath string
		exitOnError     bool
		showVersion     bool
	)

	flag.Int64Var(&intervalSec, "interval-sec", 0, "interval duration seconds for periodic command execution")
	flag.StringVar(&watchedFilePath, "watched-file", "", "file to watch; if this option is used with \"interval-sec\", it executes the command when the file was changed in the duration. If not, it executes the command when file was changed immediately")
	flag.BoolVar(&exitOnError, "exit-on-error", false, "exit this command when command execution raises the error")
	flag.BoolVar(&showVersion, "version", false, "show the version of this application")

	flag.Usage = func() {
		fmt.Fprintf(flag.CommandLine.Output(), `Usage of %s:
  %s <options> command...
options:
`, os.Args[0], os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if showVersion {
		versionJSON, _ := json.Marshal(map[string]string{
			"version":  ver,
			"revision": rev,
		})
		fmt.Printf("%s\n", versionJSON)
		return
	}

	args := flag.Args()
	if len(args) <= 0 {
		flag.Usage()
		os.Exit(1)
	}

	opt := &suez.RunOpt{
		CmdGenerator: func() *exec.Cmd {
			return exec.Command(args[0], args[1:]...)
		},
		Interval:        time.Duration(intervalSec) * time.Second,
		WatchedFilePath: watchedFilePath,
		ExitOnError:     exitOnError,
	}

	err := opt.Validate()
	if err != nil {
		flag.Usage()
		os.Exit(1)
	}

	err = suez.Run(opt)
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "%s", err)
		os.Exit(1)
	}
}
