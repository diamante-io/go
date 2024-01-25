//go:build windows
// +build windows

package ledgerbackend

import (
	"fmt"
	"os/exec"

	"github.com/Microsoft/go-winio"
)

// Windows-specific methods for the diamcircleCoreRunner type.

func (c *diamcircleCoreRunner) getPipeName() string {
	return fmt.Sprintf(`\\.\pipe\%s`, c.nonce)
}

func (c *diamcircleCoreRunner) start(cmd *exec.Cmd) (pipe, error) {
	// First set up the server pipe.
	listener, err := winio.ListenPipe(c.getPipeName(), nil)
	if err != nil {
		return pipe{}, err
	}

	// Then start the process.
	err = cmd.Start()
	if err != nil {
		listener.Close()
		return pipe{}, err
	}

	// Then accept on the server end.
	connection, err := listener.Accept()
	if err != nil {
		listener.Close()
		return pipe{}, err
	}

	return pipe{Reader: connection, File: listener}, nil
}
