package chexec

import (
	"bufio"
	"os/exec"
)

type command struct {
	*exec.Cmd

	Stdout chan []byte
	Stderr chan []byte
}

func Command(name string, args ...string) *command {
	return &command{
		Cmd:    exec.Command(name, args...),
		Stdout: make(chan []byte, 0),
		Stderr: make(chan []byte, 0),
	}
}

func (c *command) Run() error {
	stdoutReader, err := c.StdoutPipe()

	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stdoutReader)

		for scanner.Scan() {
			c.Stdout <- scanner.Bytes()
		}
	}()

	stderrReader, err := c.StderrPipe()

	if err != nil {
		return err
	}

	go func() {
		scanner := bufio.NewScanner(stderrReader)

		for scanner.Scan() {
			c.Stderr <- scanner.Bytes()
		}
	}()

	return c.Start()
}

func (c *command) Wait() chan error {
	err := make(chan error, 1)

	go func() {
		err <- c.Cmd.Wait()
	}()

	return err
}
