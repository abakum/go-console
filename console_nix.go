// +build !windows

package console

import (
	"os"
	"os/exec"

	"github.com/kr/pty"

	"github.com/runletapp/go-console/interfaces"
)

var _ interfaces.Console = (*consoleNix)(nil)

type consoleNix struct {
	file *os.File
	cmd  *exec.Cmd

	initialCols int
	initialRows int
}

func newNative(cols int, rows int) (Console, error) {
	return &consoleNix{
		initialCols: cols,
		initialRows: rows,

		file: nil,
	}, nil
}

// Start starts a process and wraps in a console
func (c *consoleNix) Start(args []string) error {
	cmd, err := c.buildCmd(args)
	if err != nil {
		return err
	}
	c.cmd = cmd

	f, err := pty.StartWithSize(cmd, &pty.Winsize{Cols: uint16(c.initialCols), Rows: uint16(c.initialRows)})
	if err != nil {
		return err
	}

	c.file = f
	return nil
}

func (c *consoleNix) buildCmd(args []string) (*exec.Cmd, error) {
	if len(args) < 1 {
		return nil, ErrInvalidCmd
	}
	cmd := exec.Command(args[0], args[1:]...)
	return cmd, nil
}

func (c *consoleNix) Read(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Read(b)
}

func (c *consoleNix) Write(b []byte) (int, error) {
	if c.file == nil {
		return 0, ErrProcessNotStarted
	}

	return c.file.Write(b)
}

func (c *consoleNix) Close() error {
	if c.file == nil {
		return ErrProcessNotStarted
	}

	return c.file.Close()
}

func (c *consoleNix) SetSize(cols int, rows int) error {
	if c.file == nil {
		c.initialRows = rows
		c.initialCols = cols
		return nil
	}

	return pty.Setsize(c.file, &pty.Winsize{Cols: uint16(cols), Rows: uint16(rows)})
}

func (c *consoleNix) GetSize() (int, int, error) {
	if c.file == nil {
		return c.initialCols, c.initialRows, nil
	}

	rows, cols, err := pty.Getsize(c.file)
	return cols, rows, err
}

func (c *consoleNix) Wait() error {
	if c.cmd == nil {
		return ErrProcessNotStarted
	}

	_, err := c.cmd.Process.Wait()
	return err
}