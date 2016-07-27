package chexec

import (
	"bufio"
	"io"
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
	var (
		err            error
		stdout, stderr io.ReadCloser
	)

	defer func() {
		if err != nil {
			c.closeChannels()
		}
	}()

	if stdout, err = c.StdoutPipe(); err != nil {
		return err
	}

	if stderr, err = c.StderrPipe(); err != nil {
		return err
	}

	err = c.Start()

	if err == nil {
		go bridge(stdout, c.Stdout)
		go bridge(stderr, c.Stderr)
	}

	return err
}

func (c *command) Wait() chan error {
	err := make(chan error, 1)

	go func() {
		err <- c.Cmd.Wait()
		c.closeChannels()
	}()

	return err
}

func (c *command) Kill() error {
	err := c.Process.Kill()
	c.closeChannels()
	return err
}

// We have potential problems from here
// How to signal to these functions that the channel(s) is(are) closed?
func (c *command) closeChannels() {
	close(c.Stdout)
	close(c.Stderr)
}

func bridge(rc io.ReadCloser, ch chan []byte) {
	scanner := bufio.NewScanner(rc)

	for scanner.Scan() {
		ch <- scanner.Bytes()
	}
}
