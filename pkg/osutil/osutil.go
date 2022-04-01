package osutil

import (
	"errors"
	"os"
	"os/exec"
	"strings"

	"github.com/rs/zerolog/log"
)

var ErrEmptyCommand = errors.New("empty command")

func RunCommand(command string) error {
	parts := strings.Split(command, " ")
	if len(parts) == 0 {
		return ErrEmptyCommand
	}

	name := parts[0]
	args := parts[1:]

	cmd := exec.Command(name, args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	log.Debug().Msg(cmd.String())

	return cmd.Run()
}
