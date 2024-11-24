//go:build !windows
// +build !windows

package main

import (
	"os"

	"github.com/abakum/cancelreader"
)

func ConsoleCP(*bool) {}

type Stdin struct {
	cancelreader.CancelReader
}

func NewStdin() (*Stdin, error) {
	cr, err := cancelreader.NewReader(os.Stdin)
	return &Stdin{
		CancelReader: cr,
	}, err
}
func (s *Stdin) Close() error {
	return s.CancelReader.Close()
}
func (s *Stdin) Cancel() bool {
	return s.CancelReader.Cancel()
}
