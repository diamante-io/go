//go:build !windows
// +build !windows

package ledgerbackend

import (
	"os"
	"os/exec"

	"github.com/pkg/errors"
)

// Posix-specific methods for the DiamcircleCoreRunner type.

func (c *diamcircleCoreRunner) getPipeName() string {
	// The exec.Cmd.ExtraFiles field carries *io.File values that are assigned
	// to child process fds counting from 3, and we'll be passing exactly one
	// fd: the write end of the anonymous pipe below.
	return "fd:3"
}

func (c *diamcircleCoreRunner) start(cmd *exec.Cmd) (pipe, error) {
	// First make an anonymous pipe.
	// Note io.File objects close-on-finalization.
	readFile, writeFile, err := os.Pipe()
	if err != nil {
		return pipe{}, errors.Wrap(err, "error making a pipe")
	}
	p := pipe{Reader: readFile, File: writeFile}

	// Add the write-end to the set of inherited file handles. This is defined
	// to be fd 3 on posix platforms.
	cmd.ExtraFiles = []*os.File{writeFile}
	err = cmd.Start()
	if err != nil {
		writeFile.Close()
		readFile.Close()
		return pipe{}, errors.Wrap(err, "error starting diamcircle-core")
	}

	return p, nil
}
