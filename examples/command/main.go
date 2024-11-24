package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/abakum/go-console"
	cc "github.com/containerd/console"
	"github.com/mattn/go-isatty"
	"github.com/xlab/closer"
)

func main() {
	var (
		raw   bool
		once  bool
		reset = func(*bool) {}
		cmd   *exec.Cmd
		arg0  = "bash"
		arg1  = "-c"
		arg2  = "echo Press any key to continue . . .;read -rn1"
	)

	defer func() {
		reset(&raw)
		closer.Close()
	}()
	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		ConsoleCP(&once)
	} else if runtime.GOOS == "windows" {
		arg0 = "cmd"
		arg1 = "/c"

		// arg0 = "powershell"
		// arg1 = "-command"

		arg2 = "pause"
	}
	log.SetFlags(log.Lmicroseconds | log.Lshortfile)
	log.SetPrefix("\r")

	var (
		// cr  cancelreader.CancelReader
		err error
		// ar  io.ReadCloser
		s *Stdin
	)
	for i := 0; i < 8; i++ {
		reset = setRaw(&raw, reset)
		if i%4 > 1 {
			reset(&raw)
			cmd = exec.Command(arg0)
			// cr, err = cancelreader.NewReader(os.Stdin)
			// if err != nil {
			// 	panic(err)
			// }
		} else {
			reset = setRaw(&raw, reset)
			cmd = exec.Command(arg0, arg1, arg2)
			// ar, err = windowsconsole.NewAnsiReaderDuplicate(os.Stdin)
			// if err != nil {
			// 	panic(err)
			// }
		}
		log.Println(cmd)
		s, err = NewStdin()
		// ar, err = windowsconsole.NewAnsiReaderDuplicate(os.Stdin)
		if err != nil {
			panic(err)
		}

		if i < 4 {
			// <Esc> <Esc> exit<Enter> exit<Enter>
			log.Println("--without PTY", i)
			cmd.Stdin = os.Stdin
			cmd.Stdout = os.Stdout
			fmt.Print("\r")
			if i%4 > 1 {
				fmt.Println("Type exit<Enter>\r")
			}
			cmd.Run()
		} else {
			// <Esc> <Esc> exit<Enter> exit<Enter>
			log.Println("--with PTY", i)
			ConsoleCP(&once)
			con, err := console.New(120, 60)
			if err != nil {
				panic(err)
			}

			fmt.Print("\r")
			if i%4 > 1 {
				fmt.Println("Type exit<Enter>\r")
			}
			err = con.Start(cmd.Args)
			if err != nil {
				panic(err)
			}

			go func() {
				_, err = io.Copy(os.Stdout, con)
				log.Println("Stdout done", i, err)
				log.Println("Cancel read stdin", i, s.Cancel())
				// if raw {
				// 	log.Println("Cancel read stdin", i, ar.Close())
				// } else {
				// 	// log.Println("Cancel read stdin", i, cr.Cancel())
				// 	log.Println("Cancel read stdin", i, ar.Close())
				// }

			}()
			_, err = io.Copy(con, s)
			// if raw {
			// 	_, err = io.Copy(con, ar)
			// } else {
			// 	// _, err = io.Copy(con, cr)
			// 	_, err = io.Copy(con, ar)
			// }
			log.Println("Stdin done", i, err)
			log.Println(con.Wait())
			log.Println("s.Close", s.Close())
			log.Println("con.Close", con.Close())
		}
	}
}

func setRaw(raw *bool, old func(*bool)) (reset func(*bool)) {
	reset = old
	if *raw {
		return
	}
	var (
		err      error
		current  cc.Console
		settings string
	)

	current, err = cc.ConsoleFromFile(os.Stdin)
	if err == nil {
		err = current.SetRaw()
		if err == nil {
			*raw = true
			reset = func(raw *bool) {
				if *raw {
					err := current.Reset()
					log.Println("Restores the console to its original state by go", err)
				}
				*raw = err != nil
			}
			log.Println("Sets the console in raw mode by go")
			return
		}
	}

	if isatty.IsCygwinTerminal(os.Stdin.Fd()) {
		settings, err = sttySettings()
		if err == nil {
			err = sttyMakeRaw()
			if err == nil {
				*raw = true
				reset = func(raw *bool) {
					if *raw {
						sttyReset(settings)
						log.Println("Restores the console to its original state by stty")
					}
					*raw = false
				}
				log.Println("Sets the console in raw mode by stty")
				return
			}
		}
	}
	log.Println(err)
	return
}

func sttyMakeRaw() error {
	cmd := exec.Command("stty", "raw", "-echo")
	cmd.Stdin = os.Stdin
	return cmd.Run()
}

func sttySettings() (string, error) {
	cmd := exec.Command("stty", "-g")
	cmd.Stdin = os.Stdin
	out, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(out)), nil
}

func sttyReset(settings string) {
	cmd := exec.Command("stty", settings)
	cmd.Stdin = os.Stdin
	_ = cmd.Run()
}
